package server

import (
	"DrawAndGuess/game"
	"DrawAndGuess/identity"
	"DrawAndGuess/user"
	"encoding/json"
	"google/uuid"
	"gorilla/mux"
	"net/http"
)

func (s *Server) listGames() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s.allGameSet.GetCurrentGames()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) listPlayersInGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameIdStr := mux.Vars(r)["gameId"]
		gameUUID, err := uuid.Parse(gameIdStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		g, err := s.allGameSet.FindCurrentGameById(gameUUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(g.PlayerSet.GetUserNames()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) getLinesInGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameIdStr := mux.Vars(r)["gameId"]
		gameUUID, err := uuid.Parse(gameIdStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		g, err := s.allGameSet.FindCurrentGameById(gameUUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(g.Lines); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) appendNewLineInGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameIdStr := mux.Vars(r)["gameId"]
		gameUUID, err := uuid.Parse(gameIdStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var lwu game.LineWithUser
		if err = json.NewDecoder(r.Body).Decode(&lwu); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		g, err := s.allGameSet.FindCurrentGameById(gameUUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		u := lwu.From
		ok, err := identity.IsIdValid(u.UserName, u.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if ok != true {
			http.Error(w, "unmatched username and id", http.StatusBadRequest)
			return
		}

		if g.DrawerName != u.UserName {
			http.Error(w, "you are not the drawer of the game", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		g.AppendLine(lwu.NewLine)
		if err = json.NewEncoder(w).Encode(g.Lines); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) setLinesInGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameIdStr := mux.Vars(r)["gameId"]
		gameUUID, err := uuid.Parse(gameIdStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var lswu game.LinesWithUser
		if err = json.NewDecoder(r.Body).Decode(&lswu); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		g, err := s.allGameSet.FindCurrentGameById(gameUUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		u := lswu.From
		ok, err := identity.IsIdValid(u.UserName, u.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if ok != true {
			http.Error(w, "unmatched username and id", http.StatusBadRequest)
			return
		}

		if g.DrawerName != u.UserName {
			http.Error(w, "you are not the drawer of the game", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		g.Lines = lswu.NewLines
		if err = json.NewEncoder(w).Encode(g.Lines); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) newGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u user.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ok, err := identity.IsIdValid(u.UserName, u.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if ok != true {
			http.Error(w, "unmatched username and id", http.StatusBadRequest)
			return
		}
		ans := mux.Vars(r)["answer"]
		user, err := s.userSet.FindUserById(u.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		g := game.NewGame(user, ans)
		user.GameId = g.Id
		if err = s.allGameSet.AppendCurrentGame(g); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(g); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) userJoinGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameIdStr := mux.Vars(r)["gameId"]
		gameUUID, err := uuid.Parse(gameIdStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var u user.User
		if err = json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		g, err := s.allGameSet.FindCurrentGameById(gameUUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		usr, err := s.userSet.FindUserById(u.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ok, err := identity.IsIdValid(usr.UserName, usr.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if ok != true {
			http.Error(w, "unmatched username and uuid", http.StatusBadRequest)
			return
		}
		if err = g.AppendPlayer(usr); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		usr.GameId = gameUUID
		if err = json.NewEncoder(w).Encode(g); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) userLeaveGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]
		gameIdStr := mux.Vars(r)["gameId"]
		userIdStr := mux.Vars(r)["userId"]

		userId, err := uuid.Parse(userIdStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ok, err := identity.IsIdValid(name, userId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if ok != true {
			http.Error(w, "unmatched name and uuid", http.StatusBadRequest)
			return
		}

		g, err := s.allGameSet.FindGameInCurrentGamesByIdStr(gameIdStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = g.DeletePlayerWithId(userId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		u, err := s.userSet.FindUserById(userId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		u.GameId = uuid.Nil

		if len(g.PlayerSet.Users) == 0 {
			s.endGame(g, "")
		}

		w.Header().Set("Content-Type", "text/plain")
		if _, err = w.Write([]byte("OK")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) listMessagesInGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameIdStr := mux.Vars(r)["gameId"]
		gameUUID, err := uuid.Parse(gameIdStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		g, err := s.allGameSet.FindCurrentGameById(gameUUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(g.Messages); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) appendMessageInGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameIdStr := mux.Vars(r)["gameId"]
		gameUUID, err := uuid.Parse(gameIdStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		g, err := s.allGameSet.FindCurrentGameById(gameUUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var m game.MessageWithUser
		err = json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		u, err := s.userSet.FindUserById(m.From.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ok, err := identity.IsIdValid(u.UserName, u.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if ok != true {
			http.Error(w, "unmatched username and id", http.StatusBadRequest)
			return
		}
		message := game.Message{
			From:    u.UserName,
			Content: m.Content,
		}
		g.Messages = append(g.Messages, message)

		if message.Content == g.Answer {
			s.endGame(g, message.From)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(g.Messages)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

func (s *Server) endGame(g *game.CurrentGame, winner string) {
	s.allGameSet.EndGame(g, winner)
}

func (s *Server) getEndedGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["gameId"]
		g, err := s.allGameSet.FindEndedGameByIdStr(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(g); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) getGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["gameId"]
		g, err := s.allGameSet.FindGameInCurrentGamesByIdStr(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "appication/json")
		if err = json.NewEncoder(w).Encode(g); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
