package migrate

import (
	"context"
	_ "embed"
	"io"
	"time"

	"github.com/jackc/pgx/v4/stdlib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
)

func systemCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "system",
		Short: "migrates the system tables of ZITADEL from one database to another",
		Long: `migrates the system tables of ZITADEL from one database to another
ZITADEL needs to be initialized
Migrations only copies keys and assets`,
		Run: func(cmd *cobra.Command, args []string) {
			config := mustNewMigrationConfig(viper.GetViper())
			copySystem(cmd.Context(), config)
		},
	}
}

func copySystem(ctx context.Context, config *Migration) {
	if instanceID == "" {
		logging.Fatal("no instance id set")
	}
	sourceClient, err := database.Connect(config.Source, false, false)
	logging.OnError(err).Fatal("unable to connect to source database")
	defer sourceClient.Close()

	destClient, err := database.Connect(config.Destination, false, true)
	logging.OnError(err).Fatal("unable to connect to destination database")
	defer destClient.Close()

	copyKeys(ctx, sourceClient, destClient)
	copyAssets(ctx, sourceClient, destClient)
}

func copyKeys(ctx context.Context, source, dest *database.DB) {
	start := time.Now()

	sourceConn, err := source.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire connection")
	defer sourceConn.Close()

	r, w := io.Pipe()
	errs := make(chan error, 1)

	go func() {
		err = sourceConn.Raw(func(driverConn interface{}) error {
			conn := driverConn.(*stdlib.Conn).Conn()
			_, err := conn.PgConn().CopyTo(ctx, w, "COPY (SELECT id, key FROM system.encryption_keys) TO stdout")
			w.Close()
			return err
		})
		errs <- err
	}()

	destConn, err := dest.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire connection")
	defer destConn.Close()

	var eventCount int64
	err = destConn.Raw(func(driverConn interface{}) error {
		conn := driverConn.(*stdlib.Conn).Conn()

		tag, err := conn.PgConn().CopyFrom(ctx, r, "COPY system.encryption_keys FROM stdin")
		eventCount = tag.RowsAffected()

		return err
	})
	logging.OnError(err).Fatal("unable to copy encryption keys to destination")
	logging.OnError(<-errs).Fatal("unable to copy encryption keys from source")
	logging.WithFields("took", time.Since(start), "count", eventCount).Info("encryption keys migrated")
}

func copyAssets(ctx context.Context, source, dest *database.DB) {
	start := time.Now()

	sourceConn, err := source.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire source connection")
	defer sourceConn.Close()

	r, w := io.Pipe()
	errs := make(chan error, 1)

	go func() {
		err = sourceConn.Raw(func(driverConn interface{}) error {
			conn := driverConn.(*stdlib.Conn).Conn()
			// ignore hash column because it's computed
			_, err := conn.PgConn().CopyTo(ctx, w, "COPY (SELECT instance_id, asset_type, resource_owner, name, content_type, data, updated_at FROM system.assets where instance_id = '"+instanceID+"') TO stdout")
			w.Close()
			return err
		})
		errs <- err
	}()

	destConn, err := dest.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire dest connection")
	defer destConn.Close()

	var eventCount int64
	err = destConn.Raw(func(driverConn interface{}) error {
		conn := driverConn.(*stdlib.Conn).Conn()

		tag, err := conn.PgConn().CopyFrom(ctx, r, "COPY system.assets FROM stdin")
		eventCount = tag.RowsAffected()

		return err
	})
	logging.OnError(err).Fatal("unable to copy assets to destination")
	logging.OnError(<-errs).Fatal("unable to copy assets from source")
	logging.WithFields("took", time.Since(start), "count", eventCount).Info("assets migrated")
}
