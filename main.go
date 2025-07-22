package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var version = "gitignore v1.0\nby PhateValleyman\nJonas.Ned@outlook.com"

var languageFiles = map[string][]string{
	"Go":            {"go.mod"},
	"Python":        {".py"},
	"Node":          {"package.json"},
	"Java":          {".java"},
	"C":             {".c"},
	"C++":           {".cpp", ".cc", ".cxx"},
	"CMake":         {"CMakeLists.txt"},
	"VisualStudio":  {".sln", ".vcxproj"},
	"Rust":          {"Cargo.toml"},
	"PHP":           {"composer.json"},
	"Swift":         {".swift"},
	"Kotlin":        {".kt", ".kts"},
}

func main() {
	helpFlag := flag.Bool("help", false, "Zobrazit napovedu")
	versionFlag := flag.Bool("version", false, "Zobrazit verzi")
	langFlag := flag.String("lang", "", "Rucni vyber jazyku (napr. --lang=Go,Python)")

	flag.BoolVar(helpFlag, "h", false, "Zobrazit napovedu")
	flag.BoolVar(versionFlag, "v", false, "Zobrazit verzi")
	flag.Usage = showHelp
	flag.Parse()

	if *helpFlag {
		showHelp()
		return
	}
	if *versionFlag {
		fmt.Println(version)
		return
	}

	checkGitRepo()

	var langs []string
	if *langFlag != "" {
		langs = strings.Split(*langFlag, ",")
	} else {
		langs = detectLanguages()
		if len(langs) == 0 {
			fmt.Println("⚠️  Nelze rozpoznat typ projektu. Používám Global.gitignore")
			langs = append(langs, "Global")
		}
	}

	confirmOverwrite(".gitignore")
	content := downloadGitignore(langs)
	os.WriteFile(".gitignore", []byte(content), 0644)

	branch := getCurrentBranch()
	gitAddCommitPush(branch, langs)

	fmt.Printf("✅ Detekováno: %v\n", langs)
}

func showHelp() {
	help := `
Použití:
  gitignore-auto [volby]

Popis:
  Automaticky rozpozná typ projektu a přidá odpovídající .gitignore z GitHubu

Volby:
  -h, --help         Zobrazí tuto nápovědu
  -v, --version      Zobrazí verzi nástroje a autora
  --lang=JAZYKY      Přepíše automatickou detekci, např. --lang=Go,Python
`
	fmt.Printf("%s\n", help)
}

func checkGitRepo() {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	if err := cmd.Run(); err != nil {
		fmt.Println("❌ Tento adresář není Git repozitář.")
		os.Exit(1)
	}
}

func detectLanguages() []string {
	detected := []string{}
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		for lang, patterns := range languageFiles {
			for _, pattern := range patterns {
				if strings.HasPrefix(pattern, ".") && strings.HasSuffix(path, pattern) {
					detected = appendIfMissing(detected, lang)
				} else if path == pattern {
					detected = appendIfMissing(detected, lang)
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("❌ Chyba při procházení souborů:", err)
	}
	return detected
}

func appendIfMissing(slice []string, item string) []string {
	for _, ele := range slice {
		if ele == item {
			return slice
		}
	}
	return append(slice, item)
}

func confirmOverwrite(file string) {
	if _, err := os.Stat(file); err == nil {
		fmt.Printf("⚠️  Soubor %s již existuje. Přepsat? (y/N): ", file)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if strings.ToLower(scanner.Text()) != "y" {
			fmt.Println("❌ Přerušeno uživatelem.")
			os.Exit(1)
		}
	}
}

func downloadGitignore(langs []string) string {
	var buffer bytes.Buffer
	for _, lang := range langs {
		url := fmt.Sprintf("https://raw.githubusercontent.com/github/gitignore/main/%s.gitignore", lang)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("❌ Chyba při stahování .gitignore:", err)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("❌ Chyba při čtení odpovědi:", err)
			continue
		}
		buffer.Write(body)
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func getCurrentBranch() string {
	cmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	branch, err := cmd.Output()
	if err != nil {
		fmt.Println("❌ Chyba při zjišťování aktuální větve:", err)
		os.Exit(1)
	}
	return strings.TrimSpace(string(branch))
}

func gitAddCommitPush(branch string, langs []string) {
	cmd := exec.Command("git", "add", ".gitignore")
	err := cmd.Run()
	if err != nil {
		fmt.Println("❌ Chyba při přidávání .gitignore do git:", err)
		os.Exit(1)
	}

	commitMessage := fmt.Sprintf("Add .gitignore for %v project", langs)
	cmd = exec.Command("git", "commit", "-m", commitMessage)
	err = cmd.Run()
	if err != nil {
		fmt.Println("❌ Chyba při commitování:", err)
		os.Exit(1)
	}

	cmd = exec.Command("git", "push", "origin", branch)
	err = cmd.Run()
	if err != nil {
		fmt.Println("❌ Chyba při pushování na origin:", err)
		os.Exit(1)
	}
}
