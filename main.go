package main

import (
	"fmt"
	"time"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

type MyFact struct {
	IntAttribute     int64
	StringAttribute  string
	BooleanAttribute bool
	FloatAttribute   float64
	TimeAttribute    time.Time
	WhatToSay        string
	Total            int64
}

func (mf *MyFact) GetWhatToSay(val int64) {
	//return fmt.Sprintf("Let say \"%s\"", sentence)
	mf.Total += val
}

func main() {
	myFact := &MyFact{
		IntAttribute:     123,
		StringAttribute:  "Some string value",
		BooleanAttribute: true,
		FloatAttribute:   1.234,
		TimeAttribute:    time.Now(),
	}

	dataCtx := ast.NewDataContext()
	err := dataCtx.Add("MF", myFact)
	if err != nil {
		panic(err)
	}

	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)

	drls := `
	rule CheckValues "Check the default values" salience 15 {
	when
		MF.IntAttribute == 123 && MF.StringAttribute == "Some string value"
	then
		 MF.GetWhatToSay(15);
		Retract("CheckValues");
	}

	rule CheckFloatAttribute "Check if FloatAttribute is greater than 1.0" salience 10 {
	when
		MF.FloatAttribute > 1.0
	then
		 MF.GetWhatToSay(10);
		Retract("CheckFloatAttribute");
	}
	`
	// Add the rule definition above into the library and name it 'TutorialRules'  version '0.0.1'

	bs := pkg.NewBytesResource([]byte(drls))
	err = ruleBuilder.BuildRuleFromResource("TutorialRules", "0.0.1", bs)
	if err != nil {
		panic(err)
	}

	knowledgeBase, _ := knowledgeLibrary.NewKnowledgeBaseInstance("TutorialRules", "0.0.1")

	engine := engine.NewGruleEngine()

	err = engine.Execute(dataCtx, knowledgeBase)
	if err != nil {
		panic(err)
	}

	fmt.Println("Response from rule engine is ", myFact.Total)
}