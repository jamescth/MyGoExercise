package parse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Nodes struct {
	BaseFields

	// The following is the json fmt for da_alarms_show output
	Node struct {
		Battery    []interface{} `json:"battery"`
		Chassis    []interface{} `json:"chassis"`
		Controller []interface{} `json:"controller"`
		Drive      []interface{} `json:"drive"`
		Metis      []interface{} `json:"metis"`
		Network    struct {
			DataPort []interface{} `json:"dataPort"`
			MgmtPort []interface{} `json:"mgmtPort"`
		} `json:"network"`
		NodeInfo struct {
			Health struct {
				Events             []interface{} `json:"events"`
				ExtendedAttributes []interface{} `json:"extendedAttributes"`
			} `json:"health"`
		} `json:"nodeInfo"`
		Pcm []interface{} `json:"pcm"`
	} `json:"node"`
	Summary []struct {
		ActiveController string `json:"activeController"`
		AmbientTemp      int    `json:"ambientTemp"`
		Beacon           string `json:"beacon"`
		Config           string `json:"config"`
		HaState          string `json:"haState"`
		HalfCapacity     bool   `json:"halfCapacity"`
		Health           struct {
			Description        string        `json:"description"`
			Events             []interface{} `json:"events"`
			ExtendedAttributes []interface{} `json:"extendedAttributes"`
			ID                 string        `json:"id"`
			Kind               string        `json:"kind"`
			Status             int           `json:"status"`
		} `json:"health"`
		Identifier     string `json:"identifier"`
		Model          string `json:"model"`
		Role           string `json:"role"`
		Serial         string `json:"serial"`
		State          string `json:"state"`
		Statetimestamp int64  `json:"statetimestamp"`
	} `json:"summary"`
}

func NewNodes(p string) (DaOut, error) {
	content, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("NewNodes"))
	}

	var n Nodes
	if err := json.Unmarshal(content, &n); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("NewNodes %s", p))
	}

	n.FileName = p
	base := path.Base(p)
	base = strings.TrimPrefix(base, "da_nodes_show.")
	base = strings.TrimRight(base, ".out")
	sTime, err := time.Parse("2006-01-02T15_04_05", base)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("NewNodes %s parse time %s", p, base))
	}

	n.StartTime = sTime.UnixNano()
	n.EndTime = sTime.UnixNano()

	var da DaOut = &n
	return da, nil

}

//func (n *Nodes) ListIssue(start, end int64, c, s string, w io.Writer) error {
func (n *Nodes) ListIssue(req *Request) (*bytes.Buffer, error) {
	var b bytes.Buffer

	// b.WriteString(fmt.Sprintf("file: %s, time: %s\n", n.Name(), daNanosecond(n.End())))
	b.WriteString(fmt.Sprintf("%20s %10s %10s %10s %10s %30s %1s\n", "ID", "Role", "State", "HA State", "Description", "State timestamp", "S"))
	b.WriteString(fmt.Sprintf("==================== ========== ========== ========== ========== ============================== =\n"))
	//fmt.Fprintf(w, "file: %s, time: %s\n", n.Name(), daNanosecond(n.End()))
	//fmt.Fprintf(w, "%20s %10s %10s %10s %10s %30s %1s\n", "ID", "Role", "State", "HA State", "Description", "State timestamp", "S")
	//fmt.Fprintf(w, "==================== ========== ========== ========== ========== ============================== =\n")
	for _, node := range n.Summary {
		//fmt.Fprintf(w, "%20s %10s %10s %10s %10s %30s %d\n",
		b.WriteString(fmt.Sprintf("%20s %10s %10s %10s %10s %30s %d\n",
			node.Health.ID,
			node.Role,
			node.State,
			node.HaState,
			node.Health.Description,
			daSecond(node.Statetimestamp),
			node.Health.Status))
	}
	return &b, nil
}
