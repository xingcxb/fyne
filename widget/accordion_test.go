package widget_test

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

func TestAccordion(t *testing.T) {
	ai := widget.NewAccordionItem("foo", widget.NewLabel("foobar"))
	t.Run("Initializer", func(t *testing.T) {
		ac := &widget.Accordion{Items: []*widget.AccordionItem{ai}}
		assert.Len(t, ac.Items, 1)
	})
	t.Run("Constructor", func(t *testing.T) {
		ac := widget.NewAccordion(ai)
		assert.Len(t, ac.Items, 1)
	})
}

func TestAccordion_Append(t *testing.T) {
	ac := widget.NewAccordion()
	ac.Append(widget.NewAccordionItem("foo", widget.NewLabel("foobar")))
	assert.Len(t, ac.Items, 1)
}

func TestAccordion_ChangeTheme(t *testing.T) {
	test.NewTempApp(t)

	ac := widget.NewAccordion()
	ac.Append(widget.NewAccordionItem("foo0", widget.NewLabel("foobar0")))
	ac.Append(widget.NewAccordionItem("foo1", widget.NewLabel("foobar1")))

	w := test.NewWindow(ac)
	defer w.Close()
	w.Resize(ac.MinSize().Add(fyne.NewSize(theme.InnerPadding(), theme.InnerPadding())))

	test.AssertImageMatches(t, "accordion/theme_initial.png", w.Canvas().Capture())

	test.WithTestTheme(t, func() {
		w.Resize(ac.MinSize().Add(fyne.NewSize(theme.InnerPadding(), theme.InnerPadding())))
		ac.Refresh()
		time.Sleep(100 * time.Millisecond)
		test.AssertImageMatches(t, "accordion/theme_changed.png", w.Canvas().Capture())
	})
}

func TestAccordion_Close(t *testing.T) {
	t.Run("Exists", func(t *testing.T) {
		ac := widget.NewAccordion()
		ac.Append(&widget.AccordionItem{
			Title:  "foo",
			Detail: widget.NewLabel("foobar"),
			Open:   true,
		})
		ac.Close(0)
		assert.False(t, ac.Items[0].Open)
		assert.False(t, ac.Items[0].Detail.Visible())
	})
	t.Run("BelowBounds", func(t *testing.T) {
		ac := widget.NewAccordion()
		ac.Append(&widget.AccordionItem{
			Title:  "foo",
			Detail: widget.NewLabel("foobar"),
			Open:   true,
		})
		ac.Close(-1)
		assert.True(t, ac.Items[0].Open)
	})
	t.Run("AboveBounds", func(t *testing.T) {
		ac := widget.NewAccordion()
		ac.Append(&widget.AccordionItem{
			Title:  "foo",
			Detail: widget.NewLabel("foobar"),
			Open:   true,
		})
		ac.Close(1)
		assert.True(t, ac.Items[0].Open)
	})
}

func TestAccordion_CloseAll(t *testing.T) {
	ac := widget.NewAccordion()
	ac.Append(widget.NewAccordionItem("foo0", widget.NewLabel("foobar0")))
	ac.Append(widget.NewAccordionItem("foo1", widget.NewLabel("foobar1")))
	ac.Append(widget.NewAccordionItem("foo2", widget.NewLabel("foobar2")))

	ac.CloseAll()
	assert.False(t, ac.Items[0].Open)
	assert.False(t, ac.Items[1].Open)
	assert.False(t, ac.Items[2].Open)
}

func TestAccordion_Layout(t *testing.T) {
	test.NewApp()

	for name, tt := range map[string]struct {
		multiOpen bool
		items     []*widget.AccordionItem
		opened    []int
	}{
		"single_open_one_item": {
			items: []*widget.AccordionItem{
				{
					Title:  "A",
					Detail: widget.NewLabel("11111"),
				},
			},
		},
		"single_open_one_item_opened": {
			items: []*widget.AccordionItem{
				{
					Title:  "A",
					Detail: widget.NewLabel("11111"),
				},
			},
			opened: []int{0},
		},
		"single_open_multiple_items": {
			items: []*widget.AccordionItem{
				{
					Title:  "A",
					Detail: widget.NewLabel("11111"),
				},
				{
					Title:  "B",
					Detail: widget.NewLabel("2222222222"),
				},
			},
		},
		"single_open_multiple_items_opened": {
			items: []*widget.AccordionItem{
				{
					Title:  "A",
					Detail: widget.NewLabel("11111"),
				},
				{
					Title:  "B",
					Detail: widget.NewLabel("2222222222"),
				},
			},
			opened: []int{0, 1},
		},
		"multiple_open_one_item": {
			multiOpen: true,
			items: []*widget.AccordionItem{
				{
					Title:  "A",
					Detail: widget.NewLabel("11111"),
				},
			},
		},
		"multiple_open_one_item_opened": {
			multiOpen: true,
			items: []*widget.AccordionItem{
				{
					Title:  "A",
					Detail: widget.NewLabel("11111"),
				},
			},
			opened: []int{0},
		},
		"multiple_open_multiple_items": {
			multiOpen: true,
			items: []*widget.AccordionItem{
				{
					Title:  "A",
					Detail: widget.NewLabel("11111"),
				},
				{
					Title:  "B",
					Detail: widget.NewLabel("2222222222"),
				},
			},
		},
		"multiple_open_multiple_items_opened": {
			multiOpen: true,
			items: []*widget.AccordionItem{
				{
					Title:  "A",
					Detail: widget.NewLabel("11111"),
				},
				{
					Title:  "B",
					Detail: widget.NewLabel("2222222222"),
				},
			},
			opened: []int{0, 1},
		},
	} {
		t.Run(name, func(t *testing.T) {
			accordion := &widget.Accordion{
				MultiOpen: tt.multiOpen,
			}
			for _, ai := range tt.items {
				accordion.Append(ai)
			}
			for _, o := range tt.opened {
				accordion.Open(o)
			}

			window := test.NewTempWindow(t, &fyne.Container{Layout: layout.NewCenterLayout(), Objects: []fyne.CanvasObject{accordion}})
			window.Resize(accordion.MinSize().Add(fyne.NewSquareSize(theme.Padding() * 2)))

			test.AssertRendersToMarkup(t, "accordion/layout_"+name+".xml", window.Canvas())
		})
	}
}

func TestAccordion_Layout_Expanded(t *testing.T) {
	test.NewApp()

	for name, tt := range map[string]struct {
		multiOpen bool
		items     []*widget.AccordionItem
		opened    []int
	}{
		"expanded_single_open_one_item": {
			items: []*widget.AccordionItem{
				{
					Title:  "A",
					Detail: widget.NewLabel("11111"),
				},
			},
		},
		"expanded_single_open_one_item_opened": {
			items: []*widget.AccordionItem{
				{
					Title:  "A",
					Detail: widget.NewLabel("11111"),
				},
			},
			opened: []int{0},
		},
		"expanded_single_open_multiple_items": {
			items: []*widget.AccordionItem{
				{
					Title:  "A",
					Detail: widget.NewLabel("11111"),
				},
				{
					Title:  "B",
					Detail: widget.NewLabel("2222222222"),
				},
			},
		},
		"expanded_single_open_multiple_items_opened": {
			items: []*widget.AccordionItem{
				{
					Title:  "A",
					Detail: widget.NewLabel("11111"),
				},
				{
					Title:  "B",
					Detail: widget.NewLabel("2222222222"),
				},
			},
			opened: []int{0, 1},
		},
		"expanded_multiple_open_one_item": {
			multiOpen: true,
			items: []*widget.AccordionItem{
				{
					Title:  "A",
					Detail: widget.NewLabel("11111"),
				},
			},
		},
		"expanded_multiple_open_one_item_opened": {
			multiOpen: true,
			items: []*widget.AccordionItem{
				{
					Title:  "A",
					Detail: widget.NewLabel("11111"),
				},
			},
			opened: []int{0},
		},
		"expanded_multiple_open_multiple_items": {
			multiOpen: true,
			items: []*widget.AccordionItem{
				{
					Title:  "A",
					Detail: widget.NewLabel("11111"),
				},
				{
					Title:  "B",
					Detail: widget.NewLabel("2222222222"),
				},
			},
		},
		"expanded_multiple_open_multiple_items_opened": {
			multiOpen: true,
			items: []*widget.AccordionItem{
				{
					Title:  "A",
					Detail: widget.NewLabel("11111"),
				},
				{
					Title:  "B",
					Detail: widget.NewLabel("2222222222"),
				},
			},
			opened: []int{0, 1},
		},
	} {
		t.Run(name, func(t *testing.T) {
			accordion := &widget.Accordion{
				MultiOpen: tt.multiOpen,
			}
			for _, ai := range tt.items {
				accordion.Append(ai)
			}
			for _, o := range tt.opened {
				accordion.Open(o)
			}

			window := test.NewTempWindow(t, &fyne.Container{Layout: layout.NewCenterLayout(), Objects: []fyne.CanvasObject{accordion}})
			window.Resize(accordion.MinSize().Max(fyne.NewSize(150, 280)))

			test.AssertRendersToMarkup(t, "accordion/layout_"+name+".xml", window.Canvas())
		})
	}
}

func TestAccordion_MinSize(t *testing.T) {
	minSizeA := fyne.MeasureText("A", theme.TextSize(), fyne.TextStyle{})
	minSizeA.Width += theme.IconInlineSize() + theme.InnerPadding()*3 + theme.Padding()
	minSizeA.Height = fyne.Max(minSizeA.Height, theme.IconInlineSize()) + theme.InnerPadding()
	minSizeB := fyne.MeasureText("B", theme.TextSize(), fyne.TextStyle{})
	minSizeB.Width += theme.IconInlineSize() + theme.InnerPadding()*3 + theme.Padding()
	minSizeB.Height = fyne.Max(minSizeB.Height, theme.IconInlineSize()) + theme.InnerPadding()
	minSize1 := fyne.MeasureText("111111", theme.TextSize(), fyne.TextStyle{})
	minSize1.Width += theme.Padding() * 4
	minSize1.Height += theme.Padding() * 4
	minSize2 := fyne.MeasureText("2222222222", theme.TextSize(), fyne.TextStyle{})
	minSize2.Width += theme.Padding() * 4
	minSize2.Height += theme.Padding() * 4

	minWidthA1 := fyne.Max(minSizeA.Width, minSize1.Width)
	minWidthB2 := fyne.Max(minSizeB.Width, minSize2.Width)
	minWidthA1B2 := fyne.Max(minWidthA1, minWidthB2)

	minHeightA1 := minSizeA.Height + minSize1.Height + theme.Padding()

	for name, tt := range map[string]struct {
		multiOpen bool
		items     []*widget.AccordionItem
		opened    []int
		want      fyne.Size
	}{
		"single_open_one_item": {
			items: []*widget.AccordionItem{
				{
					Title:  "A",
					Detail: widget.NewLabel("111111"),
				},
			},
			want: fyne.NewSize(minWidthA1, minSizeA.Height+theme.InnerPadding()),
		},
		"single_open_one_item_opened": {
			items: []*widget.AccordionItem{
				{
					Title:  "A",
					Detail: widget.NewLabel("111111"),
				},
			},
			opened: []int{0},
			want:   fyne.NewSize(minWidthA1, minHeightA1+theme.InnerPadding()),
		},
		"single_open_multiple_items": {
			items: []*widget.AccordionItem{
				{
					Title:  "A",
					Detail: widget.NewLabel("111111"),
				},
				{
					Title:  "B",
					Detail: widget.NewLabel("2222222222"),
				},
			},
			want: fyne.NewSize(minWidthA1B2, minSizeA.Height+minSizeB.Height+theme.InnerPadding()*2+theme.Padding()),
		},
		"single_open_multiple_items_opened": {
			items: []*widget.AccordionItem{
				{
					Title:  "A",
					Detail: widget.NewLabel("111111"),
				},
				{
					Title:  "B",
					Detail: widget.NewLabel("2222222222"),
				},
			},
			opened: []int{0, 1},
			want:   fyne.NewSize(minWidthA1B2, minSizeA.Height+minSizeB.Height+minSize2.Height+theme.InnerPadding()*2+theme.Padding()*2),
		},
		"multiple_open_one_item": {
			multiOpen: true,
			items: []*widget.AccordionItem{
				{
					Title:  "A",
					Detail: widget.NewLabel("111111"),
				},
			},
			want: fyne.NewSize(minWidthA1, minSizeA.Height+theme.InnerPadding()),
		},
		"multiple_open_one_item_opened": {
			multiOpen: true,
			items: []*widget.AccordionItem{
				{
					Title:  "A",
					Detail: widget.NewLabel("111111"),
				},
			},
			opened: []int{0},
			want:   fyne.NewSize(minWidthA1, minHeightA1+theme.InnerPadding()),
		},
		"multiple_open_multiple_items": {
			multiOpen: true,
			items: []*widget.AccordionItem{
				{
					Title:  "A",
					Detail: widget.NewLabel("111111"),
				},
				{
					Title:  "B",
					Detail: widget.NewLabel("2222222222"),
				},
			},
			want: fyne.NewSize(minWidthA1B2, minSizeA.Height+minSizeB.Height+theme.InnerPadding()*2+theme.Padding()),
		},
		"multiple_open_multiple_items_opened": {
			multiOpen: true,
			items: []*widget.AccordionItem{
				{
					Title:  "A",
					Detail: widget.NewLabel("111111"),
				},
				{
					Title:  "B",
					Detail: widget.NewLabel("2222222222"),
				},
			},
			opened: []int{0, 1},
			want:   fyne.NewSize(minWidthA1B2, minSizeA.Height+minSizeB.Height+minSize1.Height+minSize2.Height+theme.InnerPadding()*3+theme.Padding()),
		},
	} {
		t.Run(name, func(t *testing.T) {
			accordion := &widget.Accordion{
				MultiOpen: tt.multiOpen,
			}
			for _, ai := range tt.items {
				accordion.Append(ai)
			}
			for _, o := range tt.opened {
				accordion.Open(o)
			}

			assert.Equal(t, tt.want, accordion.MinSize())
		})
	}
}

func TestAccordion_Open(t *testing.T) {
	t.Run("Exists", func(t *testing.T) {
		ac := widget.NewAccordion()
		ac.Append(widget.NewAccordionItem("foo0", widget.NewLabel("foobar0")))
		ac.Append(widget.NewAccordionItem("foo1", widget.NewLabel("foobar1")))
		ac.Append(widget.NewAccordionItem("foo2", widget.NewLabel("foobar2")))
		assert.False(t, ac.Items[0].Open)
		assert.False(t, ac.Items[1].Open)
		assert.False(t, ac.Items[2].Open)

		ac.Open(0)
		assert.True(t, ac.Items[0].Open)
		assert.False(t, ac.Items[1].Open)
		assert.False(t, ac.Items[2].Open)

		// Opening index 1 should close index 0
		ac.Open(1)
		assert.False(t, ac.Items[0].Open)
		assert.True(t, ac.Items[1].Open)
		assert.False(t, ac.Items[2].Open)

		ac.MultiOpen = true
		ac.Open(2)
		// Opening index 2 should not close index 1
		assert.False(t, ac.Items[0].Open)
		assert.True(t, ac.Items[1].Open)
		assert.True(t, ac.Items[2].Open)
	})
	t.Run("BelowBounds", func(t *testing.T) {
		ac := widget.NewAccordion()
		ac.Append(widget.NewAccordionItem("foo", widget.NewLabel("foobar")))
		assert.False(t, ac.Items[0].Open)
		ac.Open(-1)
		assert.False(t, ac.Items[0].Open)
	})
	t.Run("AboveBounds", func(t *testing.T) {
		ac := widget.NewAccordion()
		ac.Append(widget.NewAccordionItem("foo", widget.NewLabel("foobar")))
		assert.False(t, ac.Items[0].Open)
		ac.Open(1)
		assert.False(t, ac.Items[0].Open)
	})
}

func TestAccordion_OpenAll(t *testing.T) {
	ac := widget.NewAccordion()
	ac.Append(widget.NewAccordionItem("foo0", widget.NewLabel("foobar0")))
	ac.Append(widget.NewAccordionItem("foo1", widget.NewLabel("foobar1")))
	ac.Append(widget.NewAccordionItem("foo2", widget.NewLabel("foobar2")))

	ac.OpenAll()
	// Can only open first item if !accordion.MultiOpen
	assert.True(t, ac.Items[0].Open)
	assert.False(t, ac.Items[1].Open)
	assert.False(t, ac.Items[2].Open)

	ac.MultiOpen = true
	ac.OpenAll()
	// All items should be open
	assert.True(t, ac.Items[0].Open)
	assert.True(t, ac.Items[1].Open)
	assert.True(t, ac.Items[2].Open)
}

func TestAccordion_Prepend(t *testing.T) {
	ac := widget.NewAccordion(widget.NewAccordionItem("foo", widget.NewLabel("foobar")))
	assert.Len(t, ac.Items, 1)

	ac.Prepend(widget.NewAccordionItem("bar", widget.NewLabel("more bar")))
	assert.Len(t, ac.Items, 2)
	assert.Equal(t, "bar", ac.Items[0].Title)
}

func TestAccordion_Remove(t *testing.T) {
	ai := widget.NewAccordionItem("foo", widget.NewLabel("foobar"))
	ac := widget.NewAccordion(ai)
	ac.Remove(ai)
	assert.Empty(t, ac.Items)
}

func TestAccordion_RemoveIndex(t *testing.T) {
	for name, tt := range map[string]struct {
		index  int
		length int
	}{
		"Exists":      {index: 0, length: 0},
		"BelowBounds": {index: -1, length: 1},
		"AboveBounds": {index: 1, length: 1},
	} {
		t.Run(name, func(t *testing.T) {
			ac := widget.NewAccordion()
			ac.Append(widget.NewAccordionItem("foo", widget.NewLabel("foobar")))
			ac.RemoveIndex(tt.index)
			assert.Len(t, ac.Items, tt.length)
		})
	}
}
