package main

import (
	"os"

	libcompose "github.com/docker/libcompose/yaml"
	"gopkg.in/yaml.v2"
)

// Pipeline defines the pipeline
type Pipeline struct {
	Steps yaml.MapSlice `yaml:",inline"`
}

type Clone struct {
	Git Step `yaml:"git"`
}

// DroneYAML .drone.yml
type DroneYAML struct {
	Clone    Clone         `yaml:"clone,omitempty"`
	Pipeline yaml.MapSlice `yaml:"pipeline"`
}

// When constraint
type When struct {
	Branch libcompose.Stringorslice `yaml:"branch,omitempty,flow"`
	Event  []string                 `yaml:"event,omitempty"`
	Status libcompose.Stringorslice `yaml:"status,omitempty,flow"`
}

// Step pipeline step
type Step struct {
	Image    string                   `yaml:"image"`
	Group    string                   `yaml:"group,omitempty"`
	Commands []string                 `yaml:"commands,omitempty"`
	When     When                     `yaml:"when,omitempty"`
	Vargs    map[string]interface{}   `yaml:",inline"`
	Secrets  libcompose.Stringorslice `yaml:"secrets,omitempty,flow"`
}

func specsYAML(ELBDroneHost string) []byte {
	pipeline := yaml.MapSlice{}
	downstream := yaml.MapItem{
		Key: "downstream",
		Value: Step{
			Image: "plugins/downstream",
			Vargs: map[string]interface{}{
				"server":       "http://" + ELBDroneHost,
				"repositories": []string{"bedrock/bedrock-project@${DRONE_COMMIT_BRANCH}"},
				"fork":         true,
			},
			When: When{
				Branch: []string{"dev", "qa"},
				Status: []string{"success"},
			},
			Secrets: []string{"downstream_token"},
		},
	}
	pipeline = append(pipeline, downstream)
	drone := DroneYAML{Pipeline: pipeline}
	d, _ := yaml.Marshal(drone)
	return d
}

func javaYAML(ELBDroneHost string) []byte {
	pipeline := yaml.MapSlice{}
	test := yaml.MapItem{
		Key: "test",
		Value: Step{
			Image:    "maven:3.5-jdk-8",
			Group:    "build",
			Commands: []string{"mvn clean test -B -U -s Resources/settings.xml"},
		},
	}
	downstream := yaml.MapItem{
		Key: "downstream",
		Value: Step{
			Image: "plugins/downstream",
			Vargs: map[string]interface{}{
				"server":       "http://" + ELBDroneHost,
				"repositories": []string{"bedrock/bedrock-project@${DRONE_COMMIT_BRANCH}"},
				"fork":         true,
			},
			When: When{
				Branch: []string{"dev", "qa"},
				Status: []string{"success"},
			},
			Secrets: []string{"downstream_token"},
		},
	}
	pipeline = append(pipeline, test, downstream)
	drone := DroneYAML{Pipeline: pipeline}
	d, _ := yaml.Marshal(drone)
	return d
}

func mainYAML() []byte {

	// m := make(map[string]Step)
	pipeline := yaml.MapSlice{}
	clone := Step{
		Image: "",
		Vargs: map[string]interface{}{
			"recursive":               true,
			"submodule_update_remote": true,
			"submodule_override_branch": map[string]string{
				"Java":         "${DRONE_BRANCH}",
				"ProjectSpecs": "${DRONE_BRANCH}",
			},
		},
	}
	test := yaml.MapItem{
		Key: "test",
		Value: Step{
			Image:    "maven:3-jdk-8",
			Group:    "build",
			Commands: []string{"mvn -f Java clean test -B -s Resources/settings.xml -Dmaven.repo.local=/drone/.m2"},
		},
	}

	buildFEBE := yaml.MapItem{
		Key: "build-febe",
		Value: Step{
			Image: "",
			Group: "build",
			Commands: []string{
				"febelnx project build aem --fe-parent-server  --be-parent-server http://aemas.global.febe.xms.systems --output-directory /drone/CQFiles",
				"mkdir -p /drone/CQFiles/bedrock-project/@JCR_ROOT/apps/bedrock-project/install",
				"mkdir -p /drone/CQFiles/clientlibs/@JCR_ROOT/etc/clientlibs/application/bedrock-project",
				"cp Resources/empty-folder/@nodeinfo.xml /drone/CQFiles/clientlibs/@JCR_ROOT/etc/clientlibs/application/bedrock-project",
				"cp Resources/empty-folder/@nodeinfo.xml /drone/CQFiles/bedrock-project/@JCR_ROOT/apps/bedrock-project/install"},
		},
	}

	AEM := yaml.MapItem{
		Key: "aem",
		Value: Step{
			Image: "",
			Commands: []string{
				"cd $CI_WORKSPACE",
				"mkdir -p /bin/crx-quickstart/install/install",
				"cp Resources/layerx.jcrsyncr.engine.SyncControllerImpl.config /bin/crx-quickstart/install/install/layerx.jcrsyncr.engine.SyncControllerImpl.config",
				"/bin/febe/entrypoint.sh"},
			When: When{
				Branch: []string{"dev", "qa"},
			},
			Vargs: map[string]interface{}{
				"detach": true,
			},
		},
	}

	AEMReady := yaml.MapItem{
		Key: "aem-ready",
		Value: Step{
			Image: "/drone-plugins/aem-ready",
			When: When{
				Branch: []string{"dev", "qa"},
			},
			Vargs: map[string]interface{}{
				"url": "http://aem:4502",
			},
		},
	}

	generatePackage := yaml.MapItem{
		Key: "package",
		Value: Step{
			Image: "maven:3-jdk-8",
			Commands: []string{
				"sleep 360",
				"mvn -f Java clean install -s Resources/settings.xml -B -Dcq.profilegroup=fullDeploy -Dcq.server=\"http://aem:4502\" -P-autoInstallPkgPublish -P-autoInstallPkgAuthor -Dmaven.repo.local=/drone/.m2"},
			When: When{
				Branch: []string{"dev", "qa"},
			},
		},
	}

	deployPackage := yaml.MapItem{
		Key: "deploy",
		Value: Step{
			Image: "/drone-plugins/aem-pkg-deploy",
			When: When{
				Branch: []string{"dev", "qa"},
			},
			Vargs: map[string]interface{}{
				"service_url": os.Getenv("ELB_SERVICE_DISCOVERY_HOST"),
				"company":     os.Getenv("K8_NAMESPACE"),
			},
		},
	}
	pipeline = append(pipeline, test, buildFEBE, AEM, AEMReady, generatePackage, deployPackage)
	drone := DroneYAML{Pipeline: pipeline, Clone: Clone{Git: clone}}
	d, _ := yaml.Marshal(drone)
	return d
}
