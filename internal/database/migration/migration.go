package migration

import (
	"context"

	"go.mau.fi/util/dbutil"

	"fiozap/internal/database/migration/upgrades"
	"fiozap/internal/logger"
)

const VersionTable = "fzversion"

func Run(ctx context.Context, db *dbutil.Database) error {
	logger.Component("database").Msg("running migrations")

	db.UpgradeTable = upgrades.Table

	if err := db.Upgrade(ctx); err != nil {
		return err
	}

	logger.Component("database").Msg("migrations completed")
	return nil
}
