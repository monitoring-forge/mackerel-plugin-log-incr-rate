package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/mackerelio/golib/pluginutil"
	"github.com/monitoring-forge/flagrun"
	"github.com/monitoring-forge/followparser"
	"golang.org/x/sync/errgroup"
)

var version string

type Opt struct {
	LogFile     string `long:"log-file" description:"path to log file calculate lines increased" required:"true"`
	BaseLogFile string `long:"base-log-file" description:"path to base log file count lines" required:"true"`
	KeyPrefix   string `long:"key-prefix" description:"Metric key prefix" required:"true"`
	Verbose     bool   `long:"verbose" description:"Show verbose log"`
	Version     bool   `short:"v" long:"version" description:"Show version"`
}

func (opt *Opt) run() (string, error) {
	logCounter := &simpleCounter{}
	baseLogCounter := &simpleCounter{}
	var g errgroup.Group

	g.Go(func() error {
		fp := &followparser.Parser{
			WorkDir:  pluginutil.PluginWorkDir(),
			Callback: logCounter,
			Silent:   !opt.Verbose,
		}
		_, err := fp.Parse(
			fmt.Sprintf("incr-rate-log-%s-%s", opt.KeyPrefix, url.PathEscape(opt.LogFile)),
			opt.LogFile,
		)
		if err != nil {
			return fmt.Errorf("failed to parse log file: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		fp := &followparser.Parser{
			WorkDir:  pluginutil.PluginWorkDir(),
			Callback: baseLogCounter,
			Silent:   !opt.Verbose,
		}
		_, err := fp.Parse(
			fmt.Sprintf("incr-rate-base-%s-%s", opt.KeyPrefix, url.PathEscape(opt.BaseLogFile)),
			opt.BaseLogFile,
		)
		if err != nil {
			return fmt.Errorf("failed to parse base log file: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return "", err
	}

	return opt.output(logCounter, baseLogCounter, time.Now()), nil
}

func (opt *Opt) output(logCounter, baseLogCounter *simpleCounter, now time.Time) string {
	var output strings.Builder
	timestamp := uint64(now.Unix())

	if logCounter.GetDuration() > 0 {
		fmt.Fprintf(&output, "log-incr-rate.%s_count.log\t%f\t%d\n",
			opt.KeyPrefix,
			logCounter.GetTotal()/logCounter.GetDuration(),
			timestamp)
	}
	if baseLogCounter.GetDuration() > 0 {
		fmt.Fprintf(&output, "log-incr-rate.%s_count.base\t%f\t%d\n",
			opt.KeyPrefix,
			baseLogCounter.GetTotal()/baseLogCounter.GetDuration(),
			timestamp)
	}

	if logCounter.GetDuration() > 0 && baseLogCounter.GetDuration() > 0 && baseLogCounter.GetTotal() > 0 {
		fmt.Fprintf(&output, "log-incr-rate.%s_rate.log\t%f\t%d\n",
			opt.KeyPrefix,
			(logCounter.GetTotal()/logCounter.GetDuration())/(baseLogCounter.GetTotal()/baseLogCounter.GetDuration()),
			timestamp)
	}

	return output.String()
}

func (opt *Opt) Run(_ []string) (any, int) {
	output, err := opt.run()
	if err != nil {
		return err, flagrun.CRITICAL
	}
	return output, flagrun.OK
}

func main() {
	os.Exit(flagrun.Go(&Opt{}, flagrun.Version(version)))
}
