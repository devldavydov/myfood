const RootPage = "/";

const FoodListPage = "/food";
const FoodListAPI = "/food/api/list";
const FoodGetAPI = "/food/api/get/:key";
const FoodDeleteAPI = "/food/api/del";

const NotifKey = "notifications";
const NotifClassError   = "danger";
const NotifClassWarning = "warning";
const NotifClassInfo    = "primary";

function apiPOST(url, data) {
    return $.ajax({
        url: url,
        type: "POST",
        contentType: 'application/json; charset=UTF-8',
        data: JSON.stringify(data),
        dataType: "json",
    })
}

function apiGET(url) {
    return $.ajax({
        url: url,
        type: "GET"
    })
}

function enqueueNotification(cls, msg) {    
    let nDataObj = [];

    let nDataStr = localStorage.getItem(NotifKey);
    if (nDataStr !== null)
        nDataObj = JSON.parse(nDataStr);

    nDataObj.push({
        cls: cls,
        msg: msg,
        ts: Date.now()
    })

    localStorage.setItem(NotifKey, JSON.stringify(nDataObj));
}

function showNotification(cls, msg) {
    $("#notifications").append(`
        <div class="alert alert-${cls} alert-dismissible fade show">
            ${msg}
            <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
        </div>           
    `);
}

function showPendingNotifications() {
    let nDataStr = localStorage.getItem(NotifKey);
    if (nDataStr === null)
        return;

    let nDataObj = JSON.parse(nDataStr);
    let ts = Date.now();
    for(let i=0; i < nDataObj.length; i++) {
        if (ts - nDataObj[i].ts > 10000)
            continue;
        showNotification(nDataObj[i].cls, nDataObj[i].msg);
    }

    localStorage.removeItem(NotifKey);
}

function showEl(el) {
    el.removeAttr("hidden")
}

function hideEl(el) {
    el.attr("hidden", "true")
}

// On window load.
$(function() {  
    showPendingNotifications();
});