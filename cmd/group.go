package cmd

import (
	"errors"
	"strings"
	"time"

	"github.com/Shu-AFK/bm/internal"
	"github.com/Shu-AFK/bm/internal/theme"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var GroupCmd = &cobra.Command{
	Use:   "group <add|remove|open|list>",
	Short: "Manage bookmark groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = cmd.Help()
		return nil
	},
}

var groupAddCmd = &cobra.Command{
	Use:   "add <group-name> <bm1> [bm2...]",
	Short: "Create or overwrite a group with the given bookmarks",
	RunE:  groupAdd,
}

var groupRemoveCmd = &cobra.Command{
	Use:   "remove <group-name>",
	Short: "Delete a group",
	RunE:  groupRemove,
}

var groupOpenCmd = &cobra.Command{
	Use:   "open <group-name>",
	Short: "Open all bookmarks in a group",
	RunE:  groupOpen,
}

var groupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all groups",
	RunE:  groupList,
}

func groupAdd(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		pterm.Error.Println("usage: bm group add <group-name> <bm1> [bm2...]")
		return errors.New("wrong usage")
	}

	groupName := args[0]
	bmNames := args[1:]

	bookmarks, err := internal.ReadBookmarks()
	if err != nil {
		pterm.Error.Printf("unable to read bookmarks: %v\n", err)
		return err
	}

	bmIndex := make(map[string]bool)
	for _, b := range bookmarks {
		bmIndex[b.Name] = true
	}
	for _, name := range bmNames {
		if !bmIndex[name] {
			pterm.Error.Printf("bookmark %s not found\n", name)
			return errors.New("bookmark not found")
		}
	}

	groups, err := internal.ReadGroups()
	if err != nil {
		pterm.Error.Printf("unable to read groups: %v\n", err)
		return err
	}

	for i := range groups {
		if groups[i].Name == groupName {
			groups[i].Bookmarks = bmNames
			if err := internal.WriteGroups(&groups); err != nil {
				pterm.Error.Printf("unable to write groups: %v\n", err)
				return err
			}
			pterm.Success.Printf("updated group %s\n", groupName)
			return nil
		}
	}

	groups = append(groups, internal.Group{
		Name:      groupName,
		Bookmarks: bmNames,
		CreatedAt: time.Now().Format(time.RFC3339),
	})

	if err := internal.WriteGroups(&groups); err != nil {
		pterm.Error.Printf("unable to write groups: %v\n", err)
		return err
	}

	pterm.Success.Printf("created group %s with %d bookmarks\n", groupName, len(bmNames))
	return nil
}

func groupRemove(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		pterm.Error.Println("usage: bm group remove <group-name>")
		return errors.New("wrong usage")
	}

	groupName := args[0]
	groups, err := internal.ReadGroups()
	if err != nil {
		pterm.Error.Printf("unable to read groups: %v\n", err)
		return err
	}

	found := false
	var updated []internal.Group
	for _, g := range groups {
		if g.Name == groupName {
			found = true
		} else {
			updated = append(updated, g)
		}
	}

	if !found {
		pterm.Error.Printf("group %s not found\n", groupName)
		return errors.New("group not found")
	}

	if updated == nil {
		updated = []internal.Group{}
	}

	if err := internal.WriteGroups(&updated); err != nil {
		pterm.Error.Printf("unable to write groups: %v\n", err)
		return err
	}

	pterm.Success.Printf("removed group %s\n", groupName)
	return nil
}

func groupOpen(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		pterm.Error.Println("usage: bm group open <group-name>")
		return errors.New("wrong usage")
	}

	groupName := args[0]
	groups, err := internal.ReadGroups()
	if err != nil {
		pterm.Error.Printf("unable to read groups: %v\n", err)
		return err
	}

	var group *internal.Group
	for i := range groups {
		if groups[i].Name == groupName {
			group = &groups[i]
			break
		}
	}

	if group == nil {
		pterm.Error.Printf("group %s not found\n", groupName)
		return errors.New("group not found")
	}

	bookmarks, err := internal.ReadBookmarks()
	if err != nil {
		pterm.Error.Printf("unable to read bookmarks: %v\n", err)
		return err
	}

	bmIndex := make(map[string]*internal.Bookmark)
	for i := range bookmarks {
		bmIndex[bookmarks[i].Name] = &bookmarks[i]
	}

	now := time.Now().Format(time.RFC3339)
	anyError := false
	for _, name := range group.Bookmarks {
		bm, ok := bmIndex[name]
		if !ok {
			pterm.Warning.Printf("bookmark %s not found, skipping\n", name)
			anyError = true
			continue
		}
		if err := open(bm, name); err != nil {
			anyError = true
		} else {
			bm.LastAccessed = now
			bm.AccessCount++
		}
	}

	if err := internal.WriteBookmarks(&bookmarks); err != nil {
		pterm.Warning.Printf("could not update access stats: %v\n", err)
	}

	if anyError {
		return errors.New("one or more bookmarks could not be opened")
	}
	return nil
}

func groupList(cmd *cobra.Command, args []string) error {
	groups, err := internal.ReadGroups()
	if err != nil {
		pterm.Error.Printf("unable to read groups: %v\n", err)
		return err
	}

	if len(groups) == 0 {
		pterm.Warning.Println("no groups saved yet")
		return nil
	}

	headers := []string{
		theme.GradientText("Group"),
		theme.GradientText("Bookmarks"),
	}

	data := pterm.TableData{headers}
	for _, g := range groups {
		data = append(data, []string{
			g.Name,
			strings.Join(g.Bookmarks, ", "),
		})
	}

	pterm.DefaultTable.
		WithHasHeader().
		WithSeparator("│").
		WithHeaderRowSeparator("─").
		WithBoxed(true).
		WithData(data).
		Render()

	return nil
}

func init() {
	GroupCmd.AddCommand(groupAddCmd)
	GroupCmd.AddCommand(groupRemoveCmd)
	GroupCmd.AddCommand(groupOpenCmd)
	GroupCmd.AddCommand(groupListCmd)
	rootCmd.AddCommand(GroupCmd)
}
