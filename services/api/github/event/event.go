package event

import "github.com/kevingdc/pulley/pkg/app"

type Event interface {
	Name() string
	Color() app.Color
}

type PRAssigned struct{}

func (e *PRAssigned) Name() string {
	return "Assigned"
}
func (e *PRAssigned) Color() app.Color {
	return app.ColorCyan
}

type PRUnassigned struct{}

func (e *PRUnassigned) Name() string {
	return "Unassigned"
}
func (e *PRUnassigned) Color() app.Color {
	return app.ColorGrey
}

type PRClosed struct{}

func (e *PRClosed) Name() string {
	return "Closed"
}
func (e *PRClosed) Color() app.Color {
	return app.ColorRed
}

type PRReopened struct{}

func (e *PRReopened) Name() string {
	return "Reopened"
}
func (e *PRReopened) Color() app.Color {
	return app.ColorBlue
}

type PRMerged struct{}

func (e *PRMerged) Name() string {
	return "Merged"
}
func (e *PRMerged) Color() app.Color {
	return app.ColorPurple
}

type PRReviewRequested struct{}

func (e *PRReviewRequested) Name() string {
	return "Review Requested"
}
func (e *PRReviewRequested) Color() app.Color {
	return app.ColorYellow
}

type PRApproved struct{}

func (e *PRApproved) Name() string {
	return "Approved"
}
func (e *PRApproved) Color() app.Color {
	return app.ColorGreen
}

type PRCommented struct{}

func (e *PRCommented) Name() string {
	return "Comments Added"
}
func (e *PRCommented) Color() app.Color {
	return app.ColorWhite
}

type PRChangesRequested struct{}

func (e *PRChangesRequested) Name() string {
	return "Changes Requested"
}
func (e *PRChangesRequested) Color() app.Color {
	return app.ColorOrange
}
