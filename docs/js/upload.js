// import { endpointPackage, xAuth } from 'overarching.js';
const uploadAPICall = "https://good-spec-d4rgapcc.uc.gateway.dev/package";
const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODI0NjQzNzUsIm5iZiI6MTY4MjI5MTU3NSwiaXNzIjoicGFja2l0MjMiLCJhdWQiOiJwYWNraXQyMyIsImlhdCI6MTY4MjI5MTU3NSwic3ViIjoyfQ.3ylcJEBWlwAFXiCaX1TTQ3c0INYuy_bOQMCXyOQH1hs";
// const packageData = {
//     id,
//     packagename,
//     url,
//     version,
//     author,
//     file,
//     encodedFile
// }
const formPackageName = document.getElementById("formPackageName");
const formVersionNo = document.getElementById("formVersionNo");
const formZipUpload = document.getElementById("formZipUpload");
const formURL = document.getElementById("formURL");

const errPermsMsg = document.getElementById("errPermsMsg");
const errMsg = document.getElementById("errMsg");

// var xhr = new XMLHttpRequest();
// HTML IDs
// formPackageName
// formVersionNo
// formZipUpload
// formURL
// formSubmit

// async function submitPackagebyZip(author, encodedFile, name, versionNo, zip, url) {
//     try {
//         const response = await fetch(uploadAPICall, {
//             method: 'POST',
//             headers: {
//                 xAuth: 'application/json'
//             },
//             body: JSON.stringify(data)
//         });
//         const json = await response.json();
//         console.log(json);
//     } catch (error) {
//         console.error('Error:', error);
//     }
// }

async function submitPackageByURL(inputUrl) {
    try {
        const response = await fetch(uploadAPICall, {
            method: 'POST',
            headers: {
                'X-Authorization': token
            },
            body: {
                'URL': inputUrl
            }
        });
        const json = await response.json();
        console.log(json);
    } catch (error) {
        console.error('Error:', error);
    }
    console.log('something');
}

function checkPackage() {
    // if (authenticate() == false) {
    //     errPermsMsg.style.display = "block";
    // } else 
    // if (formPackageName.value == "" || formVersionNo.value == "" || (formZipUpload.value == "" && formURL.value == "")) {
    //     errMsg.style.display = "block";
    // } else {
    //     var author = getAuthor();
    //     var encodedFile = encodeFile();
    //     submitPackage(author, encodedFile, formPackageName.value, formVersionNo.value, formZipUpload.value, formURL.value);
    // }
    // if (formZipUpload.value == "" ) {
    //     errMsg.style.display = "block";
    //     console.log("empty zip upload");
    // } 
    if (formURL.value != "") {
        submitPackageByURL(formURL.value); 
    } else if (formZipUpload.value != "" ) {
        // submitPackageByZip();
    } else {
        console.log("empty upload");
        errMsg.style.display = "block";
    }

    console.log(formZipUpload.value);
}

// function encodeFile() {

// }

// function getAuthor() {

// }