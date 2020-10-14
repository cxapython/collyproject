package mongodb

import (
"collyproject/config"
"github.com/globalsign/mgo"
	"log"
"time"
)

var session *mgo.Session

type SessionStore struct {
	session *mgo.Session
}

func init() {
	var err error
	session, err = mgo.DialWithTimeout(config.MongoConf["connect"], 5*time.Second)
	if err != nil {
		log.Panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
}

func GetS() *SessionStore {
	return &SessionStore{
		session: session.Copy(),
	}
}
//get client instance
func (s *SessionStore) GetC(cName string) *mgo.Collection {
	return s.session.DB(config.MongoConf["database"]).C(cName)
}


func (s *SessionStore) Close() {
	s.session.Close()
}