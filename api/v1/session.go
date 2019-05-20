package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-zoo/claw"

	"github.com/billyboar/battleships/helpers"
	"github.com/billyboar/battleships/models"
)

// LoadSessionRoutes will register board endpoints to /api/v1 prefix
func (s *APIServer) LoadSessionRoutes(router *mux.Router) {
	sessionRouter := router.PathPrefix("/session").Subrouter()
	c := claw.New()

	sessionRouter.HandleFunc("", s.CreateSession).Methods("POST")
	sessionRouter.Handle("", c.Use(s.GetSession).Add(s.LoadSessionToCtx)).Methods("GET")
	sessionRouter.Handle("/shoot", c.Use(s.ShootShip).Add(s.LoadSessionToCtx))
}

type SessionResponse struct {
	ID                 string              `json:"id"`
	Player             *models.Board       `json:"player"`
	ComputerShipWounds []models.Cell       `json:"computer_ship_wounds"`
	ComputerDeadShips  []models.BattleShip `json:"computer_dead_ships"`
	PlayerMissedShots  []models.Cell       `json:"player_missed_shots"`
}

func (s *APIServer) GetSession(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(SessionCtx).(*models.Session)

	response := SessionResponse{
		ID:                 session.ID,
		Player:             session.Player,
		ComputerDeadShips:  session.Computer.GetDeadShips(),
		ComputerShipWounds: session.Computer.GetAllShipWounds(),
		PlayerMissedShots:  session.Computer.MissedShots,
	}

	helpers.RenderJSON(w, response, http.StatusOK)
}

// CreateSession creates new session with randomly placed ships
// for both players
func (s *APIServer) CreateSession(w http.ResponseWriter, r *http.Request) {
	session, err := models.NewSession()
	if err != nil {
		helpers.RenderError(w, "cannot generate new session", err, http.StatusInternalServerError)
		return
	}

	event := models.CreateNewSessionEvent(session)
	if err := s.Store.AppendEvent(session.ID, event); err != nil {
		helpers.RenderError(w, "cannot append event to store", err, http.StatusInternalServerError)
		return
	}

	response := SessionResponse{
		ID:     session.ID,
		Player: session.Player,
	}

	// @TODO! return token
	helpers.RenderJSON(w, response, http.StatusCreated)
}

type ShootShipRequest struct {
	models.Cell
}

type ShootShipResponse struct {
	IsDead       bool               `json:"is_dead"`
	DeadShip     *models.BattleShip `json:"dead_ship"`
	ComputerMove struct {
		models.Cell
		DeadShip *models.BattleShip `json:"dead_ship"`
	} `json:"computer_move"`
}

// ShootShip handles shooting ships for player side
// Receives Cell as a incoming payload
func (s *APIServer) ShootShip(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var req ShootShipRequest
	if err := decoder.Decode(&req); err != nil {
		helpers.RenderError(w, "cannot decode shoot request", err, http.StatusBadRequest)
		return
	}

	if !req.IsValid() {
		helpers.RenderError(w, "shoot cell is not valid", errors.New("validation failed"), http.StatusBadRequest)
		return
	}

	session := r.Context().Value(SessionCtx).(*models.Session)

	event := models.CreateShootEvent(session.ID, &req.Cell, false)
	if err := s.Store.AppendEvent(session.ID, event); err != nil {
		helpers.RenderError(w, "cannot append event to store", err, http.StatusInternalServerError)
		return
	}

	shotStatus, deadShipID := session.Computer.RegisterShot(req.Cell)
	response := ShootShipResponse{
		IsDead: shotStatus,
	}

	if deadShip := session.Computer.MarkShipIfDead(deadShipID); deadShip != nil {
		event = models.CreateDestroyShipEvent(session.ID, deadShipID, true)
		if err := s.Store.AppendEvent(session.ID, event); err != nil {
			helpers.RenderError(w, "cannot append event to store", err, http.StatusInternalServerError)
			return
		}
		response.DeadShip = deadShip
	}

	// calculate computer response
	computerShot := session.Player.CalculateShot()
	if computerShot == nil {
		helpers.RenderError(w, "cannot find move for computer", nil, http.StatusInternalServerError)
		return
	}
	response.ComputerMove.Cell = *computerShot

	// creating shoot event for computer
	event = models.CreateShootEvent(session.ID, &response.ComputerMove.Cell, true)
	if err := s.Store.AppendEvent(session.ID, event); err != nil {
		helpers.RenderError(w, "cannot append event to store", err, http.StatusInternalServerError)
		return
	}

	response.ComputerMove.Cell.IsDead, deadShipID = session.Player.RegisterShot(response.ComputerMove.Cell)
	if deadShip := session.Player.MarkShipIfDead(deadShipID); deadShip != nil {
		event = models.CreateDestroyShipEvent(session.ID, deadShipID, false)
		if err := s.Store.AppendEvent(session.ID, event); err != nil {
			helpers.RenderError(w, "cannot append event to store", err, http.StatusInternalServerError)
			return
		}

		response.ComputerMove.DeadShip = deadShip
	}

	helpers.RenderJSON(w, response, http.StatusOK)
}
