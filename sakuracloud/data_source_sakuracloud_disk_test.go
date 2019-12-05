// Copyright 2016-2019 terraform-provider-sakuracloud authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sakuracloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccSakuraCloudDataSourceDisk_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	name := fmt.Sprintf("%s_%s", randString1, randString2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudDiskDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceDiskConfigBase(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sakuracloud_disk.disk01", "name", name),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDiskConfig(name, randString1, randString2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceExists("data.sakuracloud_disk.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "name", name),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "plan", "ssd"),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "connector", "virtio"),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "size", "20"),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "description", "source_disk_description"),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_disk.foobar", "tags.2", "tag3"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDiskConfig_With_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceExists("data.sakuracloud_disk.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDiskConfig_NotExists(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceNotExists("data.sakuracloud_disk.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceDiskConfig_With_NotExists_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceNotExists("data.sakuracloud_disk.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudDataSourceDiskConfigBase(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_disk" "disk01"{
  name = "%s"
  tags = ["tag1","tag2","tag3"]
  description = "source_disk_description"
}`, name)
}

func testAccCheckSakuraCloudDataSourceDiskConfig(name, p1, p2 string) string {
	return fmt.Sprintf(`
resource "sakuracloud_disk" "disk01"{
  name = "%s"
  tags = ["tag1","tag2","tag3"]
  description = "source_disk_description"
}

data "sakuracloud_disk" "foobar" {
  filters {
	names = ["%s", "%s"]
  }
}`, name, p1, p2)
}

func testAccCheckSakuraCloudDataSourceDiskConfig_With_Tag(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_disk" "disk01"{
  name = "%s"
  tags = ["tag1","tag2","tag3"]
  description = "source_disk_description"
}

data "sakuracloud_disk" "foobar" {
  filters {
	tags = ["tag2","tag3"]
  }
}`, name)
}

func testAccCheckSakuraCloudDataSourceDiskConfig_With_NotExists_Tag(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_disk" "disk01"{
  name = "%s"
  tags = ["tag1","tag2","tag3"]
  description = "source_disk_description"
}

data "sakuracloud_disk" "foobar" {
  filters {
	tags = ["tag2","tag4"]
  }
}`, name)
}

func testAccCheckSakuraCloudDataSourceDiskConfig_NotExists(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_disk" "disk01"{
  name = "%s"
  tags = ["tag1","tag2","tag3"]
  description = "source_disk_description"
}

data "sakuracloud_disk" "foobar" {
  filters {
	names = ["xxxxxxxxxxxxxxxxxx"]
  }
}`, name)
}
