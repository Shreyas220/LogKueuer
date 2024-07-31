package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

// Event represents a Kubernetes audit log event
type Event struct {
	Kind                     string      `json:"kind"`
	APIVersion               string      `json:"apiVersion"`
	Level                    string      `json:"level"`
	AuditID                  string      `json:"auditID"`
	Stage                    string      `json:"stage"`
	RequestURI               string      `json:"requestURI"`
	Verb                     string      `json:"verb"`
	User                     User        `json:"user"`
	SourceIPs                []string    `json:"sourceIPs"`
	UserAgent                string      `json:"userAgent"`
	ObjectRef                ObjectRef   `json:"objectRef"`
	ResponseStatus           Status      `json:"responseStatus"`
	RequestReceivedTimestamp time.Time   `json:"requestReceivedTimestamp"`
	StageTimestamp           time.Time   `json:"stageTimestamp"`
	Annotations              Annotations `json:"annotations"`
}

type User struct {
	Username string   `json:"username"`
	Groups   []string `json:"groups"`
}

type ObjectRef struct {
	Resource   string `json:"resource"`
	Namespace  string `json:"namespace"`
	Name       string `json:"name"`
	APIVersion string `json:"apiVersion"`
}

type Status struct {
	Metadata map[string]string `json:"metadata"`
	Code     int               `json:"code"`
}

type Annotations struct {
	Decision string `json:"authorization.k8s.io/decision"`
	Reason   string `json:"authorization.k8s.io/reason"`
}

func RunRandomLogsGenerator() {
	const numLogs = 1000000
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	logFile, err := os.Create("k8s_audit_logs.json")
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	defer logFile.Close()

	writer := bufio.NewWriter(logFile)

	for i := 0; i < numLogs; i++ {
		event := generateRandomEvent(rng)
		data, err := json.Marshal(event)
		if err != nil {
			log.Fatalf("Failed to marshal event: %v", err)
		}
		writer.WriteString(string(data) + "\n")
	}

	writer.Flush()
}

func generateRandomEvent(rng *rand.Rand) Event {
	return Event{
		Kind:       "Event",
		APIVersion: getRandomApi(rng),
		Level:      "Metadata",
		AuditID:    randomString(rng, 36),
		Stage:      "ResponseComplete",
		RequestURI: randomRequestURI(rng),
		Verb:       randomVerb(rng),
		User: User{
			Username: randomUsername(rng),
			Groups:   []string{"system:masters", "system:authenticated"},
		},
		SourceIPs: []string{randomIP(rng)},
		UserAgent: randomUserAgent(rng),
		ObjectRef: ObjectRef{
			Resource:   randomResources(rng),
			Namespace:  "default",
			Name:       randomPodName(rng),
			APIVersion: "v1",
		},
		ResponseStatus: Status{
			Metadata: map[string]string{},
			Code:     200,
		},
		RequestReceivedTimestamp: time.Now(),
		StageTimestamp:           time.Now().Add(time.Millisecond * time.Duration(rng.Intn(1000))),
		Annotations: Annotations{
			Decision: allowOrblock(rng),
			Reason:   "",
		},
	}
}

func getRandomApi(rng *rand.Rand) string {
	apiVersions := []string{
		"v1",
		"apps/v1",
		"batch/v1",
		"extensions/v1beta1",
		"networking.k8s.io/v1",
		"rbac.authorization.k8s.io/v1",
		"storage.k8s.io/v1",
		"apiextensions.k8s.io/v1",
		"autoscaling/v1",
		"scheduling.k8s.io/v1",
		"admissionregistration.k8s.io/v1",
		"apiregistration.k8s.io/v1",
		"certificates.k8s.io/v1",
		"coordination.k8s.io/v1",
		"discovery.k8s.io/v1",
		"node.k8s.io/v1",
		"policy/v1",
		"audit.k8s.io/v1",
		"authentication.k8s.io/v1",
		"authorization.k8s.io/v1",
		"flowcontrol.apiserver.k8s.io/v1beta1",
	}

	return apiVersions[rng.Intn(len(apiVersions))]

}

func allowOrblock(rng *rand.Rand) string {
	alow := []string{"allow", "block"}
	return alow[rng.Intn(len(alow))]
}

func randomResources(rng *rand.Rand) string {
	resources := []string{
		"pods",
		"services",
		"deployments",
		"replicasets",
		"statefulsets",
		"daemonsets",
		"jobs",
		"cronjobs",
		"configmaps",
		"secrets",
		"namespaces",
		"nodes",
		"persistentvolumes",
		"persistentvolumeclaims",
		"ingresses",
		"networkpolicies",
		"serviceaccounts",
		"roles",
		"rolebindings",
		"clusterroles",
		"clusterrolebindings",
		"resourcequotas",
		"limitranges",
		"horizontalpodautoscalers",
		"poddisruptionbudgets",
		"endpoints",
		"events",
		"replicationcontrollers",
		"leases",
		"csidrivers",
		"csinodes",
		"customresourcedefinitions",
	}

	return resources[rng.Intn(len(resources))]
}

func randomString(rng *rand.Rand, n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func randomRequestURI(rng *rand.Rand) string {
	return "/api/v1/namespaces/default/pods/" + randomPodName(rng)
}

func randomVerb(rng *rand.Rand) string {
	verbs := []string{"get", "list", "create", "update", "delete"}
	return verbs[rng.Intn(len(verbs))]
}

func randomUsername(rng *rand.Rand) string {
	usernames := []string{"kubernetes-admin", "system:serviceaccount:kube-system:default", "system:anonymous"}
	return usernames[rng.Intn(len(usernames))]
}

func randomIP(rng *rand.Rand) string {
	return fmt.Sprintf("%d.%d.%d.%d", rng.Intn(256), rng.Intn(256), rng.Intn(256), rng.Intn(256))
}

func randomUserAgent(rng *rand.Rand) string {
	agents := []string{
		"kubectl/v1.26.0 (linux/amd64) kubernetes/b46a3f8",
		"kubelet/v1.26.0 (linux/amd64) kubernetes/b46a39",
	}
	return agents[rng.Intn(len(agents))]
}

func randomPodName(rng *rand.Rand) string {
	pods := []string{"test", "nginx", "busybox", "alpine"}
	return pods[rng.Intn(len(pods))]
}
