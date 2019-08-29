// +build integration

package integration_test

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/docker/docker/pkg/ioutils"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

func assertQuickStartComponentsPresentInGit(privateSSHKeyPath string, url string) {
	dir, err := clone(privateSSHKeyPath, url)
	Expect(err).ShouldNot(HaveOccurred())
	defer os.RemoveAll(dir)
	FS := afero.Afero{Fs: afero.NewOsFs()}
	allFiles := make([]string, 0)
	err = FS.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		allFiles = append(allFiles, path)
		return nil
	})
	Expect(err).ToNot(HaveOccurred())
	fmt.Fprintf(ginkgo.GinkgoWriter, "\n all files:\n%v", allFiles)
}

type GitServer struct {
	RepoURL        string
	HostName       string
	namespace      string
	kubeconfigPath string
}

func NewGitServer(kubeconfigPath string, name string, publicSSHKeyPath string) (*GitServer, error) {
	gs := &GitServer{
		namespace:      name,
		kubeconfigPath: kubeconfigPath,
	}
	if _, err := kubectl("--kubeconfig", gs.kubeconfigPath, "create", "namespace", gs.namespace); err != nil {
		return nil, err
	}
	// Create secret to authorize public key
	createSecretArgs := []string{
		"--kubeconfig", gs.kubeconfigPath,
		"--namespace", gs.namespace,
		"create", "secret", "generic", "ssh-git",
		"--from-file=id_rsa.pub",
	}
	if _, err := kubectl(createSecretArgs...); err != nil {
		return nil, err
	}
	if _, err := kubectl("--kubeconfig", gs.kubeconfigPath, "apply", "gitsrv.yaml", gs.namespace); err != nil {
		return nil, err
	}

	// Wait for the the loadbalancer to be ready
	deadline := time.Now().Add(60 * time.Second)
	for ; time.Now().Before(deadline); time.Sleep(5 * time.Second) {
		getELBHostnameArgs := []string{
			"--namespace", gs.namespace,
			"get", "service", "gitsrv",
			"-o", "go-template='{{ (index .status.loadBalancer.ingress 0).hostname }}'",
		}
		output, err := kubectl(getELBHostnameArgs...)
		if err != nil {
			fmt.Fprintf(ginkgo.GinkgoWriter, "cannot obtain SSH Server ELB: %s, output:\n%s\nretrying ...", err, string(output))
			continue
		}
		if _, err := net.LookupIP(string(output)); err != nil {
			fmt.Fprintf(ginkgo.GinkgoWriter, "cannot resolve SSH Server ELB (%s): %s, retrying ...", string(output), err)
		}
		gs.HostName = string(output)
		gs.RepoURL = fmt.Sprint("ssh://git@%s/repos/repo.git", string(output))
		return gs, nil
	}

	return nil, errors.New("Deadline exceeded")

}

func (gs *GitServer) Delete() error {
	_, err := kubectl("--kubeconfig", gs.kubeconfigPath, "delete", "namespace", gs.namespace)
	return err
}

type SSHKeyPair struct {
	PublicKeyPath  string
	PrivateKeyPath string
}

func NewSSHKeyPair() (SSHKeyPair, error) {
	tempDir, err := ioutils.TempDir(os.TempDir(), "ssh-key-pair")
	if err != nil {
		return SSHKeyPair{}, err
	}
	sshKeyGenCmd := exec.Command("ssh-keygen", "-t", "rsa", "-N", "", "-f", "id_rsa")
	if output, err := sshKeyGenCmd.CombinedOutput(); err != nil {
		return SSHKeyPair{}, fmt.Errorf("ssh-keygen failed: %s: output:\n%s", err, string(output))
	}
	keyPair := SSHKeyPair{
		PublicKeyPath:  filepath.Join(tempDir, "id_rsa.pub"),
		PrivateKeyPath: filepath.Join(tempDir, "id_rsa"),
	}
	return keyPair, nil
}

func (kp SSHKeyPair) Delete() error {
	if err := os.Remove(kp.PrivateKeyPath); err != nil {
		return err
	}
	if err := os.Remove(kp.PublicKeyPath); err != nil {
		return err
	}
	basedir := filepath.Dir(kp.PublicKeyPath)
	return os.Remove(basedir)
}
