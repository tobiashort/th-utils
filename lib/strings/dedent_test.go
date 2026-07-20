package strings

import (
	"fmt"
)

func ExampleDedent() {
	s := `Lorem ipsum dolor sit amet,
		 |consectetur adipiscing elit.
		 |Curabitur justo tellus, facilisis nec efficitur dictum,
		 |fermentum vitae ligula. Sed eu convallis sapien.`
	fmt.Println(Dedent(s))
	// Output: Lorem ipsum dolor sit amet,
	// consectetur adipiscing elit.
	// Curabitur justo tellus, facilisis nec efficitur dictum,
	// fermentum vitae ligula. Sed eu convallis sapien.
}
