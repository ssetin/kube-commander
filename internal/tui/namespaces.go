package tui

import (
	"github.com/AnatolyRugalev/kube-commander/internal/kube"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespacesTable struct {
	*ListTable
}

func NewNamespacesTable() *NamespacesTable {
	nt := &NamespacesTable{
		ListTable: NewSelectableListTable(onNamespaceSelect),
	}
	nt.Title = "Namespaces"
	nt.resetRows()
	screen.Block.Size()
	return nt
}

func onNamespaceSelect(row []string) bool {
	screen.LoadRightPane(NewPodsTable(row[0]))
	return true
}

func (nt *NamespacesTable) resetRows() {
	nt.Rows = [][]string{
		nt.getTitleRow(),
	}
}

func (nt *NamespacesTable) getTitleRow() []string {
	return []string{"NAME", "STATUS", "AGE"}
}

func (nt *NamespacesTable) newRow(ns v1.Namespace) []string {
	return []string{
		ns.Name,
		string(ns.Status.Phase),
		Age(ns.CreationTimestamp.Time),
	}
}

func (nt *NamespacesTable) Reload() error {
	client, err := kube.GetClient()
	if err != nil {
		return err
	}
	namespaces, err := client.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	nt.resetRows()
	for _, ns := range namespaces.Items {
		nt.Rows = append(nt.Rows, nt.newRow(ns))
	}
	return nil
}