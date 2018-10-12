// Copyright 2018 Google, Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
import (
	"fmt"

	humanize "github.com/dustin/go-humanize"
)

func main() {
	var x uint64 = 1
	for i := 0; i < 15; i++ {
		fmt.Printf("%d bytes is %s\n", x, humanize.Bytes(uint64(x)))
		x = x * 10
	}
}
