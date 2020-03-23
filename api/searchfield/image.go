package searchfield

type Image string

const (
	ImageID           Image = "id"
	ImageName         Image = "name"
	ImageCommitted    Image = "committed"
	ImageDescription  Image = "description"
	ImageCreatingUser Image = "user"
)

func (i Image) String() string { return string(i) }
