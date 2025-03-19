let lastInterestString = 15;
let lastCaptalString = 2000;
let lastPeriodString = 12;

function verifyPeriods(input) {
    if (input.value > 12) {
        input.value = 12;
    } else if (input.value < 1) {
        input.value = 1;
    }

    if (input.value.length != 0) {
        lastPeriodString = input.value;
        return;
    }
    input.value = lastPeriodString;
}

function verifyInterest(input) {
    if (input.value > 100) {
        input.value = 100;
    } else if (input.value < 1) {
        input.value = 1;
    }

    if (input.value.length != 0) {
        lastInterestString = input.value;
        return;
    }
    input.value = lastInterestString;
}

function verifyAmount(input) {
    if (input.value > 20000) {
        input.value = 20000;
    } else if (input.value < 1000) {
        input.value = 1000;
    }

    if (input.value.length != 0) {
        lastCaptalString = input.value;
        return;
    }
    input.value = lastCaptalString;
}
