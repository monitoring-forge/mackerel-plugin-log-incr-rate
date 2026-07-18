package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/golib/pluginutil"
	"github.com/monitoring-forge/followparser"
	"golang.org/x/sync/errgroup"
)

var version string
var commit string

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

func main() {
	os.Exit(_main())
}

func _main() int {
	opt := &Opt{}
	psr := flags.NewParser(opt, flags.HelpFlag|flags.PassDoubleDash)
	_, err := psr.Parse()
	if opt.Version {
		if commit == "" {
			commit = "dev"
		}
		fmt.Printf(
			"%s-%s\n%s/%s, %s, %s\n",
			filepath.Base(os.Args[0]),
			version,
			runtime.GOOS,
			runtime.GOARCH,
			runtime.Version(),
			commit)
		return 0
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}

	output, err := opt.run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	fmt.Print(output)
	return 0
}
