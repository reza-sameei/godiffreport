/**
 * Copyright (c) 2019 TRIALBLAZE PTY. LTD. All rights reserved.
 *
 * Created by Reza Same'ei (reza.sameei@trialblaze.com).
 * User: Reza Same'ei
 * Date: 2019-05-27
 * Time: 08:37
 *
 * Description: diffreport_test.go
 */
package diffreport

import (
	"encoding/json"
	"gitlab.com/trialblaze/etl-common/pkg/data/formdef"
	"gitlab.com/trialblaze/etl-common/pkg/data/sdxmlchunk"
	trial "gitlab.com/trialblaze/etl-common/pkg/operations/sddm/build"
	"gitlab.com/trialblaze/etl-common/pkg/system"
	"io/ioutil"
	"testing"
)

func Debug(i interface{}) string {
	bs, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(bs)
}

func Test_DiffReport_SimpleRepr(t *testing.T) {

	past := &formdef.Representation{
		"Gender",
		map[string]string {
			"1": "M",
			"2": "Female",
		},
	}

	current := &formdef.Representation{
		"Gender",
		map[string]string{
			"1": "Male",
			"2": "Female",
			"3": "Trans",
		},
	}

	t.Log(Debug(past))

	rc := DiffReportContextNew()
	rc.Check("Form", "INIT", past, current)
	t.Logf(rc.DebugJSON())
}

func testkit_new_chunk(study string, mdv string, isetc bool, incl_study string, incl_mdv string, file string) trial.In {

	const testkit_TrialID = "1"

	bs, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	return trial.In{
		&sdxmlchunk.StudyDefinitionChunkMeta{
			testkit_TrialID,
			study,
			mdv,
			system.InternalTime{1, 1},
			"new",
			incl_study,
			incl_mdv,
			"#",
			1, // @todo fix ParsedAt type from int64 to Time
			-1,
			"at://someehere",
			nil,
			isetc,
		},
		bs,
	}
}

func testkit_build(study string, mdv string, file string) *trial.Result {

	ins := []trial.In {
		testkit_new_chunk(study, sdxmlchunk.EtcMdvOid, true, "", "", "testdata/study-etc.xml"),
		testkit_new_chunk(study, mdv, false, "", "", file),
	}

	conf := trial.BuildConf{
		"en",
	}

	rsl, _ := trial.Build(conf, ins)
	return rsl
}

func Test_DiffReport_Studies(t *testing.T) {

	a := testkit_build("S", "M", "testdata/study-at-time-001.xml")
	b := testkit_build("S", "M", "testdata/study-at-time-002.xml")

	ls := []string{}
	for k,_ := range a.Protocols {
		ls = append(ls, k)
	}
	key := ls[0]

	rc := DiffReportContextNew()

	rc.StepIn("Study", "S").StepIn("MDV", "M")

	rc.Check("Protocol", "?", a.Protocols[key], b.Protocols[key])

	t.Log(rc.DebugJSON())

}