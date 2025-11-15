package models

const MEXC = "mexc"

const FUTURES = "futures"
const SPOT = "spot"

const MIN_SPREAD = 0.07 // 7% 
const MIN_VOLUME = 0 // minimum 100_000$. TODO: not all pairs have USDT as a quote, so i need to change provided not usdt like quote to usdt