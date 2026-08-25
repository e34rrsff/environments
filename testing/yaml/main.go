package main

import (
    "os"
    "fmt"
    "go.yaml.in/yaml/v4"
    "log"
)

type Exam struct {
    Title string `yaml:"exam"`
}

type ExamList map[string]Exam

func (um *ExamList) UnmarshalYAML(value *yaml.Node) error {
    var slice []Exam
    if err := value.Decode(&slice); err != nil {
        return err
    }

    *um = make(ExamList)
    for _, exam := range slice {
        (*um)[exam.Title] = exam
    }

    return nil
}

type Config struct {
    Exams ExamList `yaml:"exams"`
}

func main() {
    data, err := os.ReadFile("exams.yaml")
    if err != nil {
        log.Fatal(err)
    }
    var cfg Config

    if err := yaml.Unmarshal(data, &cfg); err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%+v\n", cfg.Exams)
    
    os.Exit(0)
}
