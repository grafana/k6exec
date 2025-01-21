"use k6 >= 0.55.0";
import { sleep } from "k6";

export default function () {
  console.log("Hello, World!");
  sleep(0.5);
}
