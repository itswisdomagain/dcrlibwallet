package dcrlibwallet

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrd/hdkeychain"
	"github.com/decred/dcrwallet/walletseed"
)

const (
	// Approximate time (in seconds) to mine a block in mainnet
	MainNetTargetTimePerBlock = 300

	// Approximate time (in seconds) to mine a block in testnet
	TestNetTargetTimePerBlock = 120

	// Use 10% of estimated total headers fetch time to estimate rescan time
	RescanPercentage = 0.1

	// Use 80% of estimated total headers fetch time to estimate address discovery time
	DiscoveryPercentage = 0.8
)

func shutdownListener() {
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, signals...)

	// Listen for the initial shutdown signal
	select {
	case sig := <-interruptChannel:
		log.Infof("Received signal (%s).  Shutting down...", sig)
	case <-shutdownRequestChannel:
		log.Info("Shutdown requested.  Shutting down...")
	}

	// Cancel all contexts created from withShutdownCancel.
	close(shutdownSignaled)

	// Listen for any more shutdown signals and log that shutdown has already
	// been signaled.
	for {
		select {
		case <-interruptChannel:
		case <-shutdownRequestChannel:
		}
		log.Info("Shutdown signaled.  Already shutting down...")
	}
}

func contextWithShutdownCancel(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-shutdownSignaled
		cancel()
	}()
	return ctx, cancel
}

func NormalizeAddress(addr string, defaultPort string) (string, error) {
	// If the first SplitHostPort errors because of a missing port and not
	// for an invalid host, add the port.  If the second SplitHostPort
	// fails, then a port is not missing and the original error should be
	// returned.
	host, port, origErr := net.SplitHostPort(addr)
	if origErr == nil {
		return net.JoinHostPort(host, port), nil
	}
	addr = net.JoinHostPort(addr, defaultPort)
	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		return "", origErr
	}
	return addr, nil
}

// For use with gomobile bind,
// doesn't support the alternative `GenerateSeed` function because it returns more than 2 types.
func GenerateSeed() (string, error) {
	seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	if err != nil {
		return "", err
	}

	return walletseed.EncodeMnemonic(seed), nil
}

func VerifySeed(seedMnemonic string) bool {
	_, err := walletseed.DecodeUserInput(seedMnemonic)
	return err == nil
}

// ExtractDateOrTime returns the date represented by the timestamp as a date string if the timestamp is over 24 hours ago.
// Otherwise, the time alone is returned as a string.
func ExtractDateOrTime(timestamp int64) string {
	utcTime := time.Unix(timestamp, 0).UTC()
	if time.Now().UTC().Sub(utcTime).Hours() > 24 {
		return utcTime.Format("2006-01-02")
	} else {
		return utcTime.Format("15:04:05")
	}
}

func FormatUTCTime(timestamp int64) string {
	return time.Unix(timestamp, 0).UTC().Format("2006-01-02 15:04:05")
}

func AmountCoin(amount int64) float64 {
	return dcrutil.Amount(amount).ToCoin()
}

func AmountAtom(f float64) int64 {
	amount, err := dcrutil.NewAmount(f)
	if err != nil {
		log.Error(err)
		return -1
	}
	return int64(amount)
}

func EncodeHex(hexBytes []byte) string {
	return hex.EncodeToString(hexBytes)
}

func EncodeBase64(text []byte) string {
	return base64.StdEncoding.EncodeToString(text)
}

func DecodeBase64(base64Text string) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(base64Text)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func calculateTotalTimeRemaining(timeRemainingInSeconds int64) string {
	minutes := timeRemainingInSeconds / 60
	if minutes > 0 {
		return fmt.Sprintf("%d min", minutes)
	}
	return fmt.Sprintf("%d sec", timeRemainingInSeconds)
}

func calculateDaysBehind(timestamp int64) string {
	hoursBehind := float64(time.Now().Unix()-timestamp) / 60
	daysBehind := int(math.Round(hoursBehind / 24))
	if daysBehind < 1 {
		return "<1 day"
	} else if daysBehind == 1 {
		return "1 day"
	} else {
		return fmt.Sprintf("%d days", daysBehind)
	}
}

func estimateFinalBlockHeight(netType string, bestBlockTimeStamp int64, bestBlock int32) int32 {
	var targetTimePerBlock int32
	if netType == "mainnet" {
		targetTimePerBlock = MainNetTargetTimePerBlock
	} else {
		targetTimePerBlock = TestNetTargetTimePerBlock
	}

	timeDifference := time.Now().Unix() - bestBlockTimeStamp
	return (int32(timeDifference) / targetTimePerBlock) + bestBlock
}

func IsChannelClosed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}
