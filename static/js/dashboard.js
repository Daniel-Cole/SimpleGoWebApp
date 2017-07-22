var activeSearchBtn;
var activeActionBtn;
var searchType;
var actionType;
var menuWidth = 300;
var fadeInDelay = 300;
var fadeOutDelay = 250;
var appData;
var treeData;

$(document).ready(function () {

    var searchAppBtn = $("#search-app-btn");
    var searchEnvBtn = $("#search-env-btn");

    activeSearchBtn = searchAppBtn;
    searchType = "app";

    $(searchAppBtn).click(function () {
        setSelectedBtn(this, "search");
        if (actionType !== "getAll") {
            $(".control-actions-app").show();
            $(".control-actions-env").hide();
        }
        searchType = "app";
    });

    $(searchEnvBtn).click(function () {
        setSelectedBtn(this, "search");
        if (actionType !== "getAll") {
            $(".control-actions-env").show();
            $(".control-actions-app").hide();
        }
        searchType = "env";
    });

    var getBtn = $("#get-btn");
    var getAllBtn = $("#get-all-btn");
    var addBtn = $("#add-btn");
    var delBtn = $("#del-btn");

    activeActionBtn = getBtn;
    actionType = "get";

    $(getBtn).click(function () {
        setSelectedBtn(this, "action");
        actionType = "get";
        showControls();
    });

    $(getAllBtn).click(function () {
        setSelectedBtn(this, "action");
        actionType = "getAll";
        hideControls();
    });

    $(addBtn).click(function () {
        setSelectedBtn(this, "action");
        actionType = "add";
        showControls();
    });

    $(delBtn).click(function () {
        setSelectedBtn(this, "action");
        actionType = "del";
        showControls();
    });

    $("#submit-btn").click(function () {
        handleSubmit();
    });
});

function toggleNav(e) {
    var menuItems = $(document.getElementById("menu-items"));

    if (e.classList.contains("sidenav-closed")) {
        $(e).removeClass("sidenav-closed");
        $(e).addClass("sidenav-open");
        e.classList.toggle("change");
        document.getElementById("menu-div").style.width = menuWidth + "px";
        document.getElementById("main").style.marginLeft = menuWidth + "px";
        menuItems.show(fadeInDelay);
    } else {
        $(e).removeClass("sidenav-open");
        $(e).addClass("sidenav-closed");
        e.classList.toggle("change");
        document.getElementById("menu-div").style.width = "60";
        document.getElementById("main").style.marginLeft = "60";
        menuItems.hide();
    }
}

function setSelectedBtn(btn, type) {
    if (btn.classList.contains("btn-success")) {
        return;
    }
    $(btn).addClass("btn-success");
    if (type === "action") {
        if (activeActionBtn) {
            $(activeActionBtn).removeClass("btn-success");
            $(activeActionBtn).addClass("btn-default");
        }
        activeActionBtn = btn
    } else {
        if (activeSearchBtn) {
            $(activeSearchBtn).removeClass("btn-success");
            $(activeSearchBtn).addClass("btn-default");
        }
        activeSearchBtn = btn
    }
}

function hideControls() {
    $(".control-actions").fadeOut(fadeOutDelay);
}

function showControls() {
    console.log(searchType === "app");
    searchType === "app" ? $(".control-actions-app").fadeIn(fadeOutDelay) : $(".control-actions-env").fadeIn(fadeOutDelay);
}

function handleSubmit() {
    console.log('searchType:' + searchType);
    console.log('actionType:' + actionType);
    switch (searchType) {
        case "app":
            switch (actionType) {
                case "get":
                    getApplication();
                    break;
                case "getAll":
                    getAllApplications();
                    break;
                case "del":
                    deleteApplication();
                    break;
                case "add":
                    addApplication();
                    break;
            }
            break;
        case "env":
            switch (actionType) {
                case "get":
                    getEnvironment();
                    break;
                case "getAll":
                    getAllEnvironments();
                    break;
                case "del":
                    deleteEnvironment();
                    break;
                case "add":
                    addEnvironment();
                    break;
            }
            break;
    }
}

function getApplication() {
    console.log("getting application!");
    var environment = $("#env-input-1").val();
    var application = $("#app-input").val();

    console.log('env: ' + environment);
    console.log('app' + application);

    if (environment.length === 0 || application.length === 0) {
        alert("you must specify an environment and an application");
        return;
    }

    $.getJSON({
        url: '/applications/' + environment + '/' + application,
        type: 'GET',
        error: function () {
            console.log('ERROR');
            $('#info').html('<p>An error has occurred</p>');
        },
        success: function (data) {
            $('#main').html(JSON.stringify(data));
        }
    });
}

function getAllApplications() {
    $.getJSON({
        url: '/applications',
        type: 'GET',
        error: function () {
            console.log('ERROR');
            $('#info').html('<p>An error has occurred</p>');
        },
        success: function (data) {
            appData = data
            $('#main').html(JSON.stringify(data));
        }
    });
}


function getEnvironment() {
    $.getJSON({
        url: '/environment/id',
        type: 'GET',
        error: function () {
            console.log('ERROR');
            $('#info').html('<p>An error has occurred</p>');
        },
        success: function (data) {
            $('#main').html(JSON.stringify(data));
        }
    });
}

function getAllEnvironments() {
    //TODO
}

function addApplication() {

}

function deleteApplication() {

}

function addEnvironment() {

}

function deleteEnvironment() {

}

function convertAppDataToTree() {
    var flareJSON = {
        name: 'environments',
        children: []
    };


    for (application in appData) {
        //loop over each environment
        for (var i = 0; i < appData[application].length; i++) {
            app = appData[application][i];
            var primaryChild = {
                name: app.environment,
                children: []
            };

            var secondaryChild = {
                name: app.tomcat,
                children: []
            };

            primaryChild.children.push(secondaryChild);

            var exists = false;
            for (var j = 0; j < flareJSON.children.length; j++) {
                if (flareJSON.children[j].name == app.environment) {
                    //environment key already exists
                    exists = true;
                    flareJSON.children[j].children.push(secondaryChild);
                    break;
                }
            }
            if (!exists) {
                flareJSON.children.push(primaryChild);
            }

        }
    }

    treeData = flareJSON;
}



