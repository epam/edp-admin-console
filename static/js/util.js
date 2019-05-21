function getUrlParameter(sParam) {
    let sPageURL = window.location.search.substring(1),
        sURLVariables = sPageURL.split('&'),
        sParameterName,
        i;

    for (i = 0; i < sURLVariables.length; i++) {
        sParameterName = sURLVariables[i].split('=');

        if (sParameterName[0] === sParam) {
            return sParameterName[1] === undefined ? true : decodeURIComponent(sParameterName[1]);
        }
    }
}

function _sendPostRequest(async, url, data, successCallback, failCallback) {
    $.ajax({
        url: url,
        contentType: "application/json",
        type: "POST",
        data: JSON.stringify(data),
        async: async,
        success: function (resp) {
            successCallback(resp);
        },
        error: function (resp) {
            failCallback(resp);
        }
    });
}