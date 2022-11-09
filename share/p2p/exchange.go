package p2p

import (
	"bytes"
	"context"
	"github.com/celestiaorg/celestia-node/share"
	"github.com/celestiaorg/celestia-node/share/eds"
	data_request "github.com/celestiaorg/celestia-node/share/p2p/pb"
	"github.com/celestiaorg/go-libp2p-messenger/serde"
	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"io"
	"time"
)

var protocolID = protocol.ID("/eds-ex/0.0.1")
var log = logging.Logger("eds-ex")

const (
	// writeDeadline sets timeout for sending messages to the stream
	writeDeadline = time.Second * 5
	// readDeadline sets timeout for reading messages from the stream
	readDeadline = time.Minute
)

func request(ctx context.Context, to peer.ID, host host.Host, root share.Root) error {
	var dataRequest = data_request.DataRequest{
		DataRoot: root.Hash(),
	}

	stream, err := host.NewStream(ctx, to, protocolID)
	if err != nil {
		return err
	}
	if err = stream.SetWriteDeadline(time.Now().Add(writeDeadline)); err != nil {
		log.Debugf("error setting deadline: %s", err)
	}

	_, err = serde.Write(stream, &dataRequest)
	if err != nil {
		stream.Reset() //nolint:errcheck
		return err
	}
	err = stream.CloseWrite()
	if err != nil {
		log.Error(err)
	}

	if err = stream.SetReadDeadline(time.Now().Add(readDeadline)); err != nil {
		log.Debugf("error setting deadline: %s", err)
	}
	fullCar, err := io.ReadAll(stream)
	if err != nil {
		stream.Reset() //nolint:errcheck
		return err
	}

	reader := bytes.NewReader(fullCar)
	_, err = eds.ReadEDS(ctx, reader, root)
	if err != nil {
		log.Errorw("reading EDS", "err", err)
		stream.Reset() //nolint:errcheck
		return err
	}
	log.Warnf("Received EDS: %s", root.String())
	if err = stream.Close(); err != nil {
		log.Errorw("closing stream", "err", err)
	}
	return nil
}

type Server struct {
	edsStore *eds.EDSStore
}

func (s *Server) incomingRequest(stream network.Stream) {
	err := stream.SetReadDeadline(time.Now().Add(readDeadline))
	if err != nil {
		log.Debugf("error setting deadline: %s", err)
	}
	dataRequest := new(data_request.DataRequest)
	_, err = serde.Read(stream, dataRequest)
	if err != nil {
		log.Errorw("reading data request from stream", "err", err)
		stream.Reset() //nolint:errcheck
		return
	}
	if err = stream.CloseRead(); err != nil {
		log.Error(err)
	}

	ctx := context.Background()
	reader, err := s.edsStore.GetCarByBytes(ctx, dataRequest.DataRoot)
	if err != nil {
		log.Errorw("getting car by data root", "err", err)
		stream.Reset() //nolint:errcheck
		return
	}

	err = stream.SetWriteDeadline(time.Now().Add(writeDeadline))
	if err != nil {
		log.Error(err)
	}

	// TODO(distractedm1nd): This sends the entire CAR
	buf, err := io.ReadAll(reader)
	_, err = stream.Write(buf)
	if err != nil {
		log.Errorw("writing car", "err", err)
		stream.Reset() //nolint:errcheck
		return
	}

	stream.Close() //nolint:errcheck
}
