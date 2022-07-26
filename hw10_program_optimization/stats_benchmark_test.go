//nolint
package hw10programoptimization

import (
	"strings"
	"testing"
)

const (
	inputData = `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Justin Oliver Jr. Sr. I II III IV V MD DDS PhD DVM","Username":"oPerez","Email":"MelissaGutierrez@Twinte.biz","Phone":"106-05-18","Password":"f00GKr9i","Address":"Oak Valley Lane 19"}
{"Id":3,"Name":"Brian Olson","Username":"non_quia_id","Email":"FrancesEllis@Quinu.edu","Phone":"237-75-34","Password":"cmEPhX8","Address":"Butterfield Junction 74"}
{"Id":4,"Name":"Jesse Vasquez Jr. Sr. I II III IV V MD DDS PhD DVM","Username":"qRichardson","Email":"mLynch@Dabtype.name","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":5,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":6,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":7,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":8,"Name":"Jacqueline Young","Username":"CraigKnight","Email":"kCunningham@Skiptube.gov","Phone":"6-954-746-32-77","Password":"rHBCvD5JpLGs","Address":"4th Pass 91"}
{"Id":9,"Name":"Steve Burns","Username":"bRoberts","Email":"perferendis@Skippad.name","Phone":"246-85-85","Password":"68xyVtL1AaO6","Address":"Jenifer Circle 24"}
{"Id":10,"Name":"Paula Gonzales","Username":"4Ramirez","Email":"BrianBradley@Zoomcast.info","Phone":"363-62-16","Password":"atpnGIr","Address":"Barnett Park 43"}
{"Id":11,"Name":"Janet Clark I II III IV V MD DDS PhD DVM","Username":"8Perry","Email":"SaraHawkins@Eazzy.edu","Phone":"2-306-840-60-85","Password":"h7atkNURvN","Address":"Bunting Lane 40"}
{"Id":12,"Name":"Angela Nichols","Username":"quia","Email":"qui@Topiczoom.org","Phone":"9-591-154-17-64","Password":"lorrSrfSaqxk","Address":"Twin Pines Alley 66"}
{"Id":13,"Name":"Roger Gilbert","Username":"doloremque","Email":"et@Riffwire.info","Phone":"0-710-227-54-55","Password":"1lW0kUPXxj","Address":"Mifflin Lane 55"}
{"Id":14,"Name":"Donna Vasquez","Username":"aspernatur","Email":"GeraldHughes@Rhycero.name","Phone":"737-84-46","Password":"9d3cys3VNc","Address":"Lakewood Center 85"}
{"Id":15,"Name":"Todd Payne","Username":"dWilliams","Email":"3Kennedy@Plajo.gov","Phone":"0-915-255-08-38","Password":"DjkShy9NQ1R","Address":"Gulseth Parkway 72"}
{"Id":16,"Name":"Frances Olson","Username":"9Meyer","Email":"et_et@Twitterlist.biz","Phone":"0-402-112-61-97","Password":"ioOjbYTZd","Address":"Kipling Trail 73"}
{"Id":17,"Name":"Diana Palmer","Username":"quia_omnis_temporibus","Email":"AliceAustin@Realblab.biz","Phone":"288-12-53","Password":"Yyh5pLYvO7K","Address":"Butterfield Hill 19"}
{"Id":18,"Name":"Dorothy Bradley","Username":"SharonGarza","Email":"VirginiaPrice@Photospace.info","Phone":"833-89-59","Password":"TILpRLI","Address":"Judy Center 97"}
{"Id":19,"Name":"Carl Crawford","Username":"aut","Email":"dWilliams@Cogilith.info","Phone":"3-370-475-78-17","Password":"IX1rDUxz1","Address":"6th Terrace 36"}
{"Id":20,"Name":"Denise Roberts","Username":"repudiandae","Email":"quia_iusto_laboriosam@Roomm.org","Phone":"7-043-762-06-95","Password":"Ej6DnzIO5","Address":"Burrows Trail 17"}
{"Id":21,"Name":"Daniel Jenkins","Username":"PhilipFrazier","Email":"RichardPeterson@Gigashots.org","Phone":"981-77-46","Password":"JfaOx58","Address":"Morningstar Alley 24"}
{"Id":22,"Name":"Timothy Gilbert","Username":"praesentium_ut_et","Email":"et@Chatterpoint.biz","Phone":"826-94-61","Password":"VrUooZL8F8cN","Address":"Orin Plaza 40"}
{"Id":23,"Name":"Kimberly Jackson","Username":"MatthewMorales","Email":"2Snyder@Pixonyx.org","Phone":"7-452-214-79-06","Password":"wgMiOgV4","Address":"Hintze Alley 48"}
{"Id":24,"Name":"Alan Kelley","Username":"PatriciaGilbert","Email":"totam_excepturi_dolore@Oozz.biz","Phone":"9-735-960-25-03","Password":"BiDTcevghYo","Address":"Fulton Terrace 76"}
{"Id":25,"Name":"Joan Mills","Username":"in_similique_pariatur","Email":"1Webb@Rhynyx.mil","Phone":"868-93-23","Password":"uC2OdY","Address":"Bartillon Lane 57"}
{"Id":26,"Name":"Thomas Carpenter","Username":"qui_delectus_optio","Email":"jDay@Photobug.com","Phone":"4-598-198-74-05","Password":"BvKJQRFo","Address":"Northland Circle 95"}
{"Id":27,"Name":"Willie Matthews","Username":"mHenderson","Email":"autem_et_rerum@Demimbu.org","Phone":"392-26-45","Password":"0NNbpcv0oxG9","Address":"Lakeland Plaza 22"}
{"Id":28,"Name":"Rachel Andrews","Username":"yGreen","Email":"rerum@Flashpoint.mil","Phone":"7-672-375-42-30","Password":"LfEmdgN2","Address":"Melvin Hill 50"}
{"Id":29,"Name":"Laura Allen","Username":"HenryAustin","Email":"vLewis@Jabbersphere.info","Phone":"1-047-724-35-40","Password":"fyzHcHatSz4U","Address":"Mendota Parkway 20"}
{"Id":30,"Name":"Lillian Cruz","Username":"PatrickCampbell","Email":"CynthiaWashington@Pixoboo.com","Phone":"992-74-36","Password":"0pJyIl1C","Address":"Clove Hill 64"}`
)

func BenchmarkGetDomainStat(b *testing.B) {
	reader := strings.NewReader(inputData)
	for i := 0; i < b.N; i++ {
		_, err := GetDomainStat(reader, "biz")
		if err != nil {
			b.Fatalf("%v", err)
		}
		_, err = reader.Seek(0, 0)
		if err != nil {
			b.Fatalf("%v", err)
		}
	}
}
