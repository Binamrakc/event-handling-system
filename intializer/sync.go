package intializer

import "log"

func DBmigrate() {
	err := DB.AutoMigrate(&User{}, &Auth{}, &ContactMessage{}, &Event{})
	if err != nil {
		log.Println("auto migration failed", err)
		return
	}

}
