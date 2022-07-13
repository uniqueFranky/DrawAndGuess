package server

import (
	"DrawAndGuess/identity"
	"DrawAndGuess/storage"
	"DrawAndGuess/user"
	"encoding/json"
	"google/uuid"
	"gorilla/mux"
	"net/http"
)

func (s *Server) newUserRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]
		psw := mux.Vars(r)["psw"]
		_, err := s.userSet.UserReg(name, psw)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", " application/json")
		u := user.User{
			UserName: name,
			UserId:   uuid.Nil,
		}
		if err = json.NewEncoder(w).Encode(u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) newUserLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]
		psw := mux.Vars(r)["psw"]
		id, err := s.userSet.UserLogin(name, psw)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		u, err := s.userSet.FindUserById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) listOnlineUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s.userSet.GetUserNames()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) userLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]
		idStr := mux.Vars(r)["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		u, err := s.userSet.FindUserById(id)
		ok, err := identity.IsIdValid(name, id)
		if false == ok || err != nil {
			http.Error(w, "unmatched id and username", http.StatusBadRequest)
			return
		}
		err = s.userSet.DeleteUserById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		g, err := s.allGameSet.FindGameInCurrentGamesById(u.GameId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = g.DeletePlayerWithId(u.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) getUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["name"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		u, err := s.userSet.FindUserById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func (s *Server) hasReg() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]
		rows, err := storage.NewQuery("select name from users where name = '" + name + "';")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		if rows.Next() {
			_, err = w.Write([]byte("Y"))
		} else {
			_, err = w.Write([]byte("N"))
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
