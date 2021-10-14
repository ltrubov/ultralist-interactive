package cmd

func Map(src []string, f func(string) string) []string {
  res := make([]string, len(src))
  for i, v := range src {
      res[i] = f(v)
  }
  return res
}

func Contains(src []string, target string) bool {
  for _,el := range src {
    if el == target {
      return true
    }
  }

  return false
}

// func RemoveHelpArgs(src []string) []string {
//   var res = make([]string, len(src))
//   var rc = 0
//   for i,v := range src {
//     if strings.ToLower(comps[0]) == "uhelp" ||
//        Contains(comps, "-h") ||
//        Contains(comps, "--help") {
//         rc++
//         continue
//     }
//     res[i] = v
//   }
//
//   return res[:len(src)-rc]
// }