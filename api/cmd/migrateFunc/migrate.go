package migrateFunc

import (
	"fmt"
	"strings"
	"suffgo/cmd/config"
	"suffgo/cmd/database"
	o "suffgo/internal/options/infrastructure/models"
	p "suffgo/internal/proposals/infrastructure/models"
	r "suffgo/internal/rooms/infrastructure/models"
	s "suffgo/internal/settingsRoom/infrastructure/models"
	ur "suffgo/internal/userRooms/infrastructure/models"
	m "suffgo/internal/users/infrastructure/models"
	e "suffgo/internal/votes/infrastructure/models"
)

func Make() error {
	conf := config.GetConfig()
	db := database.NewPostgresDatabase(conf)

	err := MigrateUser(db)
	
	if err != nil {
		return err
	}

	err = MigrateRoom(db)
	if err != nil {
		return err
	}

	err = MigrateProposal(db) 

	if err != nil {
		return err
	}

	
	err = MigrateOption(db)

	if err != nil {
		return err
	}

	err = MigrateVote(db)
	if err != nil {
		return err
	}


	err = MigrateRoomSetting(db)
	if err != nil {
		return err
	}

	err = MakeConstraints(db)
	if err != nil {
		fmt.Printf("Error al agregar la clave foránea: %v\n", err)
		panic(err.Error())
	}

	return nil
}

func MigrateUser(db database.Database) error {
	err := db.GetDb().Sync2(new(m.Users))

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Se ha migrado User con exito\n")
	}

	return err
}

func MigrateRoom(db database.Database) error {
	err := db.GetDb().Sync2(new(r.Room), new(ur.UserRoom))

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Se ha migrado Room y UserRoom con exito\n")
	}

	return err
}

func MigrateProposal(db database.Database) error {
	err := db.GetDb().Sync2(new(p.Proposal))

	if err != nil {
		return err
	} else {
		fmt.Printf("Se ha migrado Proposal con exito\n")
	}

	return nil
}

func MigrateOption(db database.Database) error {
	err := db.GetDb().Sync2(new(o.Option))

	if err != nil {
		return err
	} else {
		fmt.Printf("Se ha migrado Option con exito\n")
	}
	return nil
}

func MigrateRoomSetting(db database.Database) error {
	err := db.GetDb().Sync2(new(s.SettingsRoom))

	if err != nil {
		return err
	} else {
		fmt.Printf("Se ha migrado RoomSetting con exito\n")
	}

	return nil
}

func MigrateVote(db database.Database) error {
	err := db.GetDb().Sync2(new(e.Vote))

	if err != nil {
		return err
	} else {
		fmt.Printf("Se ha migrado Election con exito\n")
	}

	return nil
}

func MakeConstraints(db database.Database) error {
    statements := []struct {
        sql  string
        info string
    }{
        {
            `ALTER TABLE user_room ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)`,
            "fk_user on user_room",
        },
        {
            `ALTER TABLE user_room ADD CONSTRAINT fk_room FOREIGN KEY (room_id) REFERENCES room(id)`,
            "fk_room on user_room",
        },
        {
            `ALTER TABLE room DROP CONSTRAINT IF EXISTS fk_user; ALTER TABLE room ADD CONSTRAINT fk_user FOREIGN KEY (admin_id) REFERENCES users(id)`,
            "fk_user on room",
        },
        {
            `ALTER TABLE proposal ADD CONSTRAINT fk_room FOREIGN KEY (room_id) REFERENCES room(id)`,
            "fk_room on proposal",
        },
        {
            `ALTER TABLE settings_room ADD CONSTRAINT fk_room FOREIGN KEY (room_id) REFERENCES room(id)`,
            "fk_room on settings_room",
        },
        {
            `ALTER TABLE option ADD CONSTRAINT fk_proposal FOREIGN KEY (proposal_id) REFERENCES proposal(id) ON DELETE CASCADE`,
            "fk_proposal with ON DELETE CASCADE",
        },
        {
            `ALTER TABLE vote ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)`,
            "fk_user on vote",
        },
        {
            `ALTER TABLE vote ADD CONSTRAINT fk_option FOREIGN KEY (option_id) REFERENCES option(id)`,
            "fk_option on vote",
        },
        {
            `CREATE UNIQUE INDEX IF NOT EXISTS value_proposal_idx ON option(value, proposal_id)`,
            "value_proposal_idx unique index on option(value, proposal_id)",
        },
    }

    for _, stmt := range statements {
        if err := execIgnoreExists(db, stmt.sql); err != nil {
            return fmt.Errorf("error ejecutando '%s': %w", stmt.info, err)
        }
        fmt.Printf("%s creada o ya existía, OK\n", stmt.info)
    }

    return nil
}


// execIgnoreExists ejecuta la sentencia SQL y, si falla porque la constraint
// ya existe, devuelve nil. Para otros errores, los propaga.
func execIgnoreExists(db database.Database, sql string) error {
    if _, err := db.GetDb().Exec(sql); err != nil {
        if strings.Contains(err.Error(), "already exists") {
            // Constraint ya existe: lo ignoramos
            return nil
        }
        return err
    }
    return nil
}