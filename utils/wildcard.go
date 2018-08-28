package utils

//static int wildcmp(const char* wild, const char* pattern)
//{
//  // Written by Jack Handy - &lt;A href="mailto:jakkhandy@hotmail.com"&gt;jakkhandy@hotmail.com&lt;/A&gt;
//  const char *cp = NULL, *mp = NULL;

//  while ((*pattern) && (*wild != '*'))
//  {
//    if ((*wild != *pattern) && (*wild != '?'))
//    {
//      return 0;
//    }
//    wild++;
//    pattern++;
//  }

//  while (*pattern)
//  {
//    if (*wild == '*')
//    {
//      if (!*++wild)
//      {
//        return 1;
//      }
//      mp = wild;
//      cp = pattern + 1;
//    }
//    else if ((*wild == *pattern) || (*wild == '?'))
//    {
//      wild++;
//      pattern++;
//    }
//    else
//    {
//      wild = mp;
//      pattern = cp++;
//    }
//  }

//  while (*wild == '*')
//  {
//    wild++;
//  }
//  return !*wild;
//}

//// The main function that checks if two given strings match. The first
//// string may contain wildcard characters
//static bool match(const char *first,  const char* second)
//{
//  // If we reach at the end of both strings, we are done
//  if (*first == '\0' && *second == '\0')
//    return true;

//  // Make sure that the characters after '*' are present in second string.
//  // This function assumes that the first string will not contain two
//  // consecutive '*'
//  if (*first == '*' && *(first + 1) != '\0' && *second == '\0')
//    return false;

//  // If the first string contains '?', or current characters of both
//  // strings match
//  if (*first == '?' || *first == *second)
//    return match(first + 1, second + 1);

//  // If there is *, then there are two possibilities
//  // a) We consider current character of second string
//  // b) We ignore current character of second string.
//  if (*first == '*')
//    return match(first + 1, second) || match(first, second + 1);
//  return false;
//}

//wirdcard match, see https://en.wikipedia.org/wiki/Wild_card
func WildcardCmp(txt, pattern string) bool {
	txtLen := len(txt)
	patternLen := len(pattern)

	if patternLen == 0 {
		return false
	}

	var txtIdx int
	var patternIdx int
	var mark int

	for txtIdx < txtLen && patternIdx < patternLen {
		if pattern[patternIdx] == '?' {
			patternIdx++
			txtIdx++
			continue
		} else if pattern[patternIdx] == '*' {
			patternIdx++
			mark = patternIdx
			continue
		}

		if pattern[patternIdx] != txt[txtIdx] {
			if patternIdx == 0 && txtIdx == 0 {
				return false
			}
			txtIdx -= (patternIdx - mark - 1)
			patternIdx = mark

			continue
		}

		txtIdx++
		patternIdx++
	}

	if patternIdx == patternLen {
		if txtIdx == txtLen {
			return true
		}
		if pattern[patternIdx-1] == byte('*') {
			return true
		}
	}

	for patternIdx < patternLen {
		if pattern[patternIdx] != byte('*') {
			return false
		}
		patternIdx++
	}
	return txtIdx == txtLen
}
