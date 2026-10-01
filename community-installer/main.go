package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func steamLibraries(steam string) ([]string, error) {
	roots := []string{steam}
	data, exists, err := readFile(filepath.Join(steam, "steamapps", "libraryfolders.vdf"))
	if err != nil {
		return nil, err
	}
	if !exists {
		return roots, nil
	}
	nodes, err := parseVDF(string(data))
	if err != nil {
		return nil, err
	}
	root, err := vdfChild(nodes, "libraryfolders")
	if err != nil {
		return nil, err
	}
	if root == nil || root.scalar {
		return nil, fmt.Errorf("invalid Steam library list")
	}
	for _, n := range root.children {
		if n.scalar {
			continue
		}
		path, err := vdfChild(n.children, "path")
		if err != nil {
			return nil, err
		}
		if path == nil || !path.scalar {
			continue
		}
		absolutePath, err := absolute(path.value)
		if err != nil {
			return nil, err
		}
		duplicate := false
		for _, known := range roots {
			if samePath(known, absolutePath) {
				duplicate = true
			}
		}
		if !duplicate {
			roots = append(roots, absolutePath)
		}
	}
	return roots, nil
}

func readLine(reader *bufio.Reader, label string) (string, error) {
	fmt.Print(label)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("input cancelled")
	}
	return strings.Trim(strings.TrimSpace(line), `"`), nil
}
func choose(reader *bufio.Reader, label string, values []string, interactive bool) (string, error) {
	if len(values) == 1 {
		return values[0], nil
	}
	if !interactive {
		return "", fmt.Errorf("%s: supply an explicit path/account (found %d candidates)", label, len(values))
	}
	if len(values) == 0 {
		return readLine(reader, label+" path: ")
	}
	fmt.Println(label + ":")
	for i, path := range values {
		fmt.Printf("  %d. %s\n", i+1, path)
	}
	answer, err := readLine(reader, "Selection number: ")
	if err != nil {
		return "", err
	}
	number, err := strconv.Atoi(answer)
	if err != nil || number < 1 || number > len(values) {
		return "", fmt.Errorf("invalid selection")
	}
	return values[number-1], nil
}

func resolveLocations(game, profile, steam, account, config string, reader *bufio.Reader, interactive bool) (locations, error) {
	var loc locations
	var err error
	if steam == "" {
		steam = registrySteam()
	}
	if steam == "" {
		candidate := filepath.Join(os.Getenv("ProgramFiles(x86)"), "Steam")
		if _, e := os.Stat(candidate); e == nil {
			steam = candidate
		}
	}
	if steam == "" && (game == "" || config == "") {
		steam, err = choose(reader, "Steam installation", nil, interactive)
		if err != nil {
			return loc, err
		}
	}
	if steam != "" {
		steam, err = absolute(steam)
		if err != nil {
			return loc, err
		}
	}
	if game == "" {
		libraries, err := steamLibraries(steam)
		if err != nil {
			return loc, err
		}
		var games []string
		for _, library := range libraries {
			candidate := filepath.Join(library, "steamapps", "common", "ProjectZomboid")
			if _, e := os.Stat(filepath.Join(candidate, "ProjectZomboid64.json")); e == nil {
				games = append(games, candidate)
			}
		}
		game, err = choose(reader, "Project Zomboid installation", games, interactive)
		if err != nil {
			return loc, err
		}
	}
	loc.Game, err = absolute(game)
	if err != nil {
		return loc, err
	}
	if config == "" {
		var configs []string
		if account != "" {
			if strings.Trim(account, "0123456789") != "" {
				return loc, fmt.Errorf("Steam account directory must be numeric")
			}
			configs = []string{filepath.Join(steam, "userdata", account, "config", "localconfig.vdf")}
		} else {
			configs, err = filepath.Glob(filepath.Join(steam, "userdata", "*", "config", "localconfig.vdf"))
			if err != nil {
				return loc, err
			}
		}
		config, err = choose(reader, "Steam account configuration", configs, interactive)
		if err != nil {
			return loc, err
		}
	}
	loc.SteamConfig, err = absolute(config)
	if err != nil {
		return loc, err
	}
	if profile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return loc, err
		}
		profile = filepath.Join(home, "Zomboid")
		data, exists, err := readFile(loc.SteamConfig)
		if err != nil || !exists {
			return loc, fmt.Errorf("cannot read Steam account")
		}
		opts, _, err := lookupVDF(data, optionsPath)
		if err != nil {
			return loc, err
		}
		tokens, err := launchTokens(opts)
		if err != nil {
			return loc, err
		}
		for _, t := range tokens {
			if strings.HasPrefix(t.value, "-cachedir=") {
				profile = strings.TrimPrefix(t.value, "-cachedir=")
				if !filepath.IsAbs(profile) {
					return loc, fmt.Errorf("relative -cachedir requires separate review")
				}
			}
		}
	}
	loc.Profile, err = absolute(profile)
	return loc, err
}

func showPlan(p plan) {
	fmt.Printf("\nRuntime: %s (verified embedded package and Community signature)\n", runtimeVersion)
	fmt.Printf("Game: %s\nProfile: %s\nSteam account: %s\n", p.Locations.Game, p.Locations.Profile, p.Locations.SteamConfig)
	fmt.Println("Normal Windows Steam launcher only. No game launch or subscription changes.")
	fmt.Println("One native bootstrap loads the Community Java framework; existing mod approvals and saves are retained.")
	fmt.Println("Local mods will take priority over Workshop copies, including other duplicate Mod IDs.")
	fmt.Println("Keep the original Workshop dependency downloaded if another mod needs it.")
	if len(p.Changes) == 0 {
		fmt.Println("Already configured: all managed files and launch settings match.")
		return
	}
	fmt.Printf("Changes (%d files; backups will be stored in the game directory):\n", len(p.Changes))
	for _, c := range p.Changes {
		path, _ := targetPath(p.Locations, c.Area, c.Name, p.Payload)
		action := "Create"
		if c.BeforeExists {
			action = "Replace"
		}
		if c.Area == "steam" {
			action = "Update only PZ LaunchOptions in"
		}
		fmt.Printf("  %s %s\n", action, path)
	}
}

func run() error {
	fmt.Printf("ZombieBuddy Community Installer %s\nSource: %s\n", installerVersion, sourceCommit)
	if runtime.GOOS != "windows" {
		return fmt.Errorf("this preview installer supports Windows only")
	}
	game := flag.String("game", "", "Project Zomboid game directory")
	profile := flag.String("profile", "", "Zomboid cache/profile directory (default: user home or Steam -cachedir)")
	steam := flag.String("steam", "", "Steam installation directory")
	account := flag.String("account", "", "numeric Steam userdata directory")
	config := flag.String("steam-config", "", "exact selected account localconfig.vdf")
	preview := flag.Bool("plan", false, "read-only preview; no files are changed")
	install := flag.Bool("install", false, "install after showing the plan")
	rollbackPath := flag.String("rollback", "", "restore the installation recorded by receipt.json")
	yes := flag.Bool("yes", false, "accept the printed plan (process and integrity checks still apply)")
	flag.Parse()
	if flag.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments")
	}
	if *preview && (*install || *rollbackPath != "") {
		return fmt.Errorf("choose one operation")
	}
	if *install && *rollbackPath != "" {
		return fmt.Errorf("choose one operation")
	}
	payload, err := payloadFiles(embeddedPackage)
	if err != nil {
		return err
	}
	interactive := !*preview && !*install && *rollbackPath == ""
	reader := bufio.NewReader(os.Stdin)
	if interactive {
		answer, err := readLine(reader, "[I] Install / [R] Roll back / [Q] Quit: ")
		if err != nil {
			return err
		}
		switch strings.ToLower(answer) {
		case "i", "install":
		case "r", "rollback":
			value, err := readLine(reader, "Path to receipt.json: ")
			if err != nil {
				return err
			}
			*rollbackPath = value
		case "q", "quit", "":
			return nil
		default:
			return fmt.Errorf("unknown operation")
		}
	}
	if *rollbackPath != "" {
		path, err := absolute(*rollbackPath)
		if err != nil {
			return err
		}
		fmt.Printf("Roll back managed files using %s\n", path)
		if !*yes {
			answer, err := readLine(reader, "Continue? [y/N]: ")
			if err != nil {
				return err
			}
			if !strings.EqualFold(answer, "y") {
				fmt.Println("Cancelled. No changes applied.")
				return nil
			}
		}
		if err := rollback(path, payload, stoppedGuard); err != nil {
			return err
		}
		fmt.Println("Rollback complete. Original managed files/settings restored; unrelated Steam changes retained.")
		return nil
	}
	loc, err := resolveLocations(*game, *profile, *steam, *account, *config, reader, interactive)
	if err != nil {
		return err
	}
	p, err := buildPlan(loc, payload, gameHash)
	if err != nil {
		return err
	}
	showPlan(p)
	if *preview || len(p.Changes) == 0 {
		return nil
	}
	if !*yes {
		answer, err := readLine(reader, "Close Steam and the game. Install now? [y/N]: ")
		if err != nil {
			return err
		}
		if !strings.EqualFold(answer, "y") {
			fmt.Println("Cancelled. No changes applied.")
			return nil
		}
	}
	receiptPath, err := applyPlan(p, stoppedGuard, nil)
	if err != nil {
		return err
	}
	fmt.Printf("Installation verified. Receipt and rollback backups: %s\n", receiptPath)
	fmt.Println("Start Steam normally, enable ZombieBuddy Community for your save, and keep dependent mods unchanged.")
	fmt.Println("This development preview still requires game acceptance; no gameplay/MP compatibility claim is made.")
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
	if len(os.Args) == 1 {
		fmt.Print("Press Enter to close...")
		bufio.NewReader(os.Stdin).ReadString('\n')
	}
	if err != nil {
		os.Exit(1)
	}
}
