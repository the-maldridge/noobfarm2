package netauth

import (
	"context"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/netauth/netauth/pkg/netauth"
	"github.com/the-maldridge/noobfarm2/internal/web"
	"github.com/the-maldridge/noobfarm2/internal/web/auth"
)

func init() {
	auth.RegisterCallback(cb)
}

func cb() {
	auth.RegisterFactory("netauth", New)
}

// New obtains a new authentication service that uses the NetAuth
// backend.
func New(l hclog.Logger) (web.Auth, error) {
	l = l.Named("netauth")
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/netauth/")
	viper.AddConfigPath("$HOME/.netauth")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		l.Error("Fatal error reading configuration", "error", err)
		return nil, err
	}

	// Grab a client
	c, err := netauth.New()
	if err != nil {
		l.Error("Error during NetAuth initialization", "error", err)
		return nil, err
	}
	c.SetServiceName("terrastate")

	group := os.Getenv("NF_AUTH_GROUP")

	x := netAuthBackend{
		nacl:  c,
		l:     l,
		group: group,
	}

	return &x, nil
}

type netAuthBackend struct {
	nacl *netauth.Client
	l    hclog.Logger

	group string
}

func (b *netAuthBackend) AuthUser(ctx context.Context, user, pass string) error {
	err := b.nacl.AuthEntity(ctx, user, pass)
	if status.Code(err) != codes.OK {
		return err
	}

	groups, err := b.nacl.EntityGroups(ctx, user)
	if status.Code(err) != codes.OK {
		b.l.Warn("RPC Error: ", "error", err)
		return err
	}

	for _, g := range groups {
		b.l.Trace("Checking group for user", "user", user, "want", b.group, "have", g.GetName())
		if g.GetName() == b.group {
			b.l.Debug("User authenticated", "group", group, "user", user)
			return nil
		}
	}

	return auth.ErrUnauthenticated
}
