$("#readloanpublic").on("click", function () {
    const loanuid = $("#loanuid").val();
    window.location = "http://localhost:3000/readloan/public/" + loanuid;
});

$("#readloanprivate").on("click", function () {
    const loanuid = $("#loanuid").val();
    window.location = "http://localhost:3000/readloan/private/" + loanuid;
});