// console.log("HIIII")

const main = "main-list"

// https://www.w3schools.com/howto/howto_js_collapsible.asp
// may be helpful later
/*

window.onload = () => {
    console.log("LOADED")
    fetch("http://localhost:9999/requests/entry/")
    .then(async (response) => {
        console.log("ATTEMPTING")

        if (response.ok) {
            var dat = (await response.json());
            console.log(dat)
        }

    }).catch((err) => {
        console.log("Error! " + err)
    })

}


interface ArticleData {
	url: string;
	title:         string;    
	fromWeb:       string;   
	summary:       string;    
	timeAdded:     string;
	timePublished: string;
}*/

const entriesPerPage = 5;
var entries = 0; 
var mainList;

async function GetData(page) {
    var ret = await fetch("http://localhost:9999/requests/entry")
    .then(async (response) => {

        if (response.ok) {
            var dat = (await response.json());
            console.log(dat);
            return dat;
        }

    }).catch((err) => {
        console.log("Error! " + err)
    })

    return ret;

}


function AddDisplayDat(doc, data) {

    if (doc === null) {
        console.log("COULDNT FIND DOC");
        return; 
    }


    const newDat = document.createElement("entry" + entries);
    newDat.textContent = data.summary;

    doc.appendChild(newDat);
    
    
    
}



window.onload = async function () {
    mainList = document.getElementById("main-list");

    
    var dat_ = await GetData(0);
    AddDisplayDat(mainList, dat_);

  };