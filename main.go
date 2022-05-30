package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const TMP_FOLDER_LOCATION = "tmp_git_file/"

// Program that downloads 1 file from a git repository
//

// git-file <link repo>
// git-file https://github.com/Yadiiiig/ydb/blob/master/internals/proto/insert.pb.go

// Later-on add ssh profiles

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Not enough arguments")
        os.Exit(0)
	} else if len(args) < 3 {
		args = append(args, ".")
	}

	err := download(args[1], args[2])
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Successfully downloaded file.")
}

func download(url, location string) error {
	object := parse(url)
	path := fmt.Sprintf("%s%s/", TMP_FOLDER_LOCATION, object.RepoName)

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	_, err = command(fmt.Sprintf("git clone %s %s", object.Url, path))
	if err != nil {
		return err
	}

	_, err = command(fmt.Sprintf("mv %s/%s %s", path, object.FilePath, location))
	if err != nil {
		return err
	}

	_, err = command(fmt.Sprintf("rm -rf -d %s", path))
	if err != nil {
		return err
	}
	return nil
}

type object struct {
	Url      string
	RepoName string
	FilePath string
}

func parse(value string) object {
	url := value
	if strings.HasPrefix(url, "https://") {
		url = strings.TrimPrefix(url, "https://")
	}

	split := strings.Split(url, "/blob/")
	path := strings.Split(split[1], "/")[1:]

	return object{
		Url:      strings.Split(value, "/blob/")[0],
		RepoName: strings.Split(split[0], "/")[len(split)],
		FilePath: strings.Join(path, "/"),
	}
}

func command(c string) (string, error) {
	out, err := exec.Command("sh", "-c", c).Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
