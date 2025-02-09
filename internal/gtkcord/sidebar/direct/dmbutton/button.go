package dmbutton

import (
	"context"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotkit/gtkutil"
	"github.com/diamondburned/gotkit/gtkutil/cssutil"
	"github.com/diamondburned/gtkcord4/internal/gtkcord"
	"github.com/diamondburned/gtkcord4/internal/gtkcord/sidebar/sidebutton"
	"github.com/diamondburned/ningen/v3"
	"github.com/diamondburned/ningen/v3/states/read"
)

type Button struct {
	*gtk.Overlay
	Pill   *sidebutton.Pill
	Button *gtk.Button

	ctx context.Context
}

var dmButtonCSS = cssutil.Applier("sidebar-dm-button-overlay", `
	.sidebar-dm-button {
		padding: 4px 12px;
		border-radius: 0;
	}
`)

func NewButton(ctx context.Context, open func()) *Button {
	b := Button{ctx: ctx}

	icon := gtk.NewImageFromIconName("user-available")
	icon.SetIconSize(gtk.IconSizeLarge)
	icon.SetPixelSize(gtkcord.GuildIconSize)

	b.Button = gtk.NewButton()
	b.Button.AddCSSClass("sidebar-dm-button")
	b.Button.SetTooltipText("Direct Messages")
	b.Button.SetChild(icon)
	b.Button.SetHasFrame(false)
	b.Button.SetHAlign(gtk.AlignCenter)
	b.Button.ConnectClicked(func() {
		b.Pill.State = sidebutton.PillActive
		b.Pill.Invalidate()

		open()
	})

	b.Pill = sidebutton.NewPill()

	b.Overlay = gtk.NewOverlay()
	b.Overlay.SetChild(b.Button)
	b.Overlay.AddOverlay(b.Pill)

	vis := gtkutil.WithVisibility(ctx, b)

	state := gtkcord.FromContext(ctx)
	state.BindHandler(vis, func(ev gateway.Event) {
		switch ev := ev.(type) {
		case *read.UpdateEvent:
			if ev.GuildID.IsValid() {
				return
			}

			b.Invalidate()
		}
	})

	dmButtonCSS(b)
	return &b
}

// Invalidate forces a complete recheck of all direct messaging channels to
// update the unread indicator.
func (b *Button) Invalidate() {
	state := gtkcord.FromContext(b.ctx)
	unread := dmUnreadState(state)

	b.Pill.Attrs = sidebutton.PillAttrsFromUnread(unread)
	b.Pill.Invalidate()
}

func dmUnreadState(state *gtkcord.State) ningen.UnreadIndication {
	var unread ningen.UnreadIndication

	chs, _ := state.Cabinet.PrivateChannels()
	for _, ch := range chs {
		state := state.ChannelIsUnread(ch.ID)
		if state > unread {
			unread = state
		}

		if unread == ningen.ChannelMentioned {
			break
		}
	}

	return unread
}
