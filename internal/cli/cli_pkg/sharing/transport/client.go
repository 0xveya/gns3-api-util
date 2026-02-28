package transport

import (
	"context"
	"crypto/ed25519"
	"crypto/tls"
	"errors"
	"fmt"

	"github.com/0xveya/gns3util/internal/cli/cli_pkg/sharing/keys"
	"github.com/0xveya/gns3util/internal/cli/cli_pkg/sharing/trust"
	"github.com/0xveya/gns3util/internal/cli/cli_pkg/utils/colorUtils"
	"github.com/quic-go/quic-go"
)

type VerifyPrompt func(peerLabel, fp string, words []string) (bool, error)

func DialWithPin(
	ctx context.Context,
	addr string,
	myLabel string,
	myFP string,
	myPub ed25519.PublicKey,
	ts *trust.Store,
	prompt VerifyPrompt,
) (*quic.Conn, *quic.Stream, Hello, error) {
	tlsConf := &tls.Config{
		NextProtos:         []string{"gns3util/1"},
		InsecureSkipVerify: true, // #nosec G402
		MinVersion:         tls.VersionTLS13,
	}

	conn, err := quic.DialAddr(ctx, addr, tlsConf, &quic.Config{})
	if err != nil {
		return nil, nil, Hello{}, err
	}

	ctrl, openErr := conn.OpenStreamSync(ctx)
	if openErr != nil {
		_ = conn.CloseWithError(0, "open stream failed")
		return nil, nil, Hello{}, openErr
	}

	// 1) Hello
	if writeHelloErr := WriteJSON(ctx, ctrl, Hello{Label: myLabel, FP: myFP}); writeHelloErr != nil {
		_ = conn.CloseWithError(0, "hello failed")
		return nil, nil, Hello{}, writeHelloErr
	}
	var srv Hello
	if jsonErr := ReadJSON(ctx, ctrl, &srv); jsonErr != nil {
		_ = conn.CloseWithError(0, "hello recv failed")
		return nil, nil, Hello{}, jsonErr
	}

	// 2) SAS nonce exchange
	clientNonce, nonceErr := NewNonce()
	if nonceErr != nil {
		_ = conn.CloseWithError(0, "nonce gen failed")
		return nil, nil, Hello{}, nonceErr
	}
	if writeNonceErr := WriteJSON(ctx, ctrl, SASMsg{Nonce: clientNonce}); writeNonceErr != nil {
		_ = conn.CloseWithError(0, "sas write failed")
		return nil, nil, Hello{}, writeNonceErr
	}
	var srvSAS SASMsg
	if readJsonErr := ReadJSON(ctx, ctrl, &srvSAS); readJsonErr != nil {
		_ = conn.CloseWithError(0, "sas read failed")
		return nil, nil, Hello{}, readJsonErr
	}

	// 3) Extract server public key from TLS to compute fingerprint and SAS
	st := conn.ConnectionState().TLS
	if len(st.PeerCertificates) == 0 {
		_ = conn.CloseWithError(0, "no peer cert")
		return nil, nil, Hello{}, errors.New("no peer certificate")
	}
	cert := st.PeerCertificates[0]
	serverPub, ok := cert.PublicKey.(ed25519.PublicKey)
	if !ok {
		_ = conn.CloseWithError(0, "unexpected key type")
		return nil, nil, Hello{}, errors.New("server key is not ed25519")
	}
	serverFP := keys.Fingerprint(serverPub)

	// 4) Derive SAS code bound to server identity + fresh nonces
	words, deriveErr := DerivePGPWordsSimple(serverPub, clientNonce, srvSAS.Nonce, 3)
	if deriveErr != nil {
		_ = conn.CloseWithError(0, "sas derive failed")
		return nil, nil, Hello{}, deriveErr
	}
	fmt.Printf("%s %s\n", colorUtils.Info("Verify code:"), colorUtils.Highlight(FormatSAS(words)))

	// 5) Pinning: if not pinned, ask the user to accept
	if _, ok := ts.Get(serverFP); !ok {
		if prompt == nil {
			_ = conn.CloseWithError(0, "unpinned and no prompt")
			return nil, nil, Hello{}, errors.New("unpinned and no prompt")
		}
		accept, promptErr := prompt(srv.Label, serverFP, words)
		if promptErr != nil || !accept {
			_ = conn.CloseWithError(0, "verification rejected")
			if promptErr == nil {
				promptErr = errors.New("verification rejected")
			}
			return nil, nil, Hello{}, promptErr
		}
		// Persist pin
		_ = ts.Add(serverFP, srv.Label)
	}

	return conn, ctrl, srv, nil
}
