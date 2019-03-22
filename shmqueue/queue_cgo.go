// +build !cgo

package shmqueue

/*

#define SMQ_OK -1
#define SMQ_QUIT -2
#define SMQ_FULL -3

type shmQueue struct{
	unsigned char* pBuff;
	unsigned long long* pRIdx;
	unsigned long long* pWIdx;
	unsgined long long each;
	unsgined long long total;
	unsigned long pQ;
}SHM_QUEUE;

extern int sqSemTake(SHM_QUEUE* pQ);

extern int sqSemGive(SHM_QUEUE* pQ, block int);

int sqBPush(SHM_QUEUE* pQ, const void* p, unsigned long size) {
	unsigned long long remain = *pQ->pWIdx - *pQ->pRIdx;
	if (remain == pQ->total) {
		return SMQ_FULL;
	}
	unsigned long long offset = remain % pQ->total;
	int valueOffset = int(offset * pQ->each);
	
	if (size > (unsigned long)pQ->each) {
		size = pQ->each - 4;
	}	
	memcpy(pQ->pBuff+valueOffset, (void*)&size, 4);
	memcpy(pQ->pBuff+valueOffset+4, p, size);
	if (__sync_fetch_and_add(pQ->pWIdx, 1) == 1) {
		sqSemGive(pQ);
	}
	return size;
}

int sqBPop(SHM_QUEUE* pQ, void* p, unsigned long size) {
	for (;;) {
		unsigned long long remain = *pQ->pWIdx - *pQ->pRIdx;		
		if (remain == 0) {
			int status = sqSemTake(pQ);
			if (status == SMQ_QUIT) {
				return SMQ_QUIT;
			}
		} else {
			unsigned long long offset = remain % pQ->total;
			int valueOffset = int(offset * pQ->each);
			unsigned long realSize = 0;
			memcpy((void*)&size, pQ->pBuff+valueOffset, 4);
			if (size < realSize) {
				size = realSize
			}
			memcpy(p, pQ->pBuff+valueOffset+4, size);	
			__sync_fetch_and_add(pQ->pRIdx, 1);
			return size;
		}
	}
}
*/
import "C"
import "unsafe"

//export sqSemTake
func sqSemTake(p *C.SHM_QUEUE) C.int {
	q := (*shmQueue)unsafe.Pointer(p.pQ)
	sem := q.sem
	
	if	sem.Take() {
		return C.SMQ_QUIT
	}
	return C.SMQ_OK	
}

//export sqSemGive
func sqSemGive(p *C.SHM_QUEUE, block int) C.int {
	q := (*shmQueue)unsafe.Pointer(p.pQ)
	sem := q.sem
	if block==1 {
		sem.Give(true)
	} else {
		sem.Give(false)
	}
	return SMQ_OK
}

