

interface ArticleData {
	url: string;
	title:         string;    
	fromWeb:       string;   
	summary:       string;    
	timeAdded:     string;
	timePublished: string;
}

const entriesPerPage = 5;
var entries = 0; 
var mainList: HTMLElement; 

async function GetData(page: number): Promise<ArticleData | null> {
    var ret = fetch("http://localhost:9999/requests/entry")
    .then(async (response) => {

        if (response.ok) {
            var dat: ArticleData = (await response.json()) as ArticleData;
            return dat;
        }

    }).catch((err) => {
        console.log("Error! " + err)
    })

    return ret as Promise<ArticleData>;

}


function AddDisplayDat(doc: HTMLElement | null,  data: ArticleData) {

    if (doc === null) {
        console.log("no doc")
        return; 
    }


    const newDat = document.createElement("entry" + entries);
    newDat.textContent = data.summary;

    doc.appendChild(newDat);
    
    
}


// function AddCollap



window.onload = async function () {
    mainList = document.getElementById("main-list");
    /*await GetData(0).then((dat: ArticleData) => {
        AddDisplayDat(mainList, dat);
    });*/
  };


