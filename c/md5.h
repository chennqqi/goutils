#ifndef __MD5C_H__
#define __MD5C_H__

/* POINTER defines a generic pointer type */
typedef unsigned char * POINTER;


/* UINT4 defines a four byte word */
/*typedef unsigned long int UINT32; */


/* MD5 context. */
typedef struct {
	unsigned int state[4];                                   /* state (ABCD) */
	unsigned int count[2];        /* number of bits, modulo 2^64 (lsb first) */
	unsigned char buffer[64];                         /* input buffer */
} MD5_CTX;


#ifdef __cplusplus
extern "C" {
#endif

	void MD5Init(MD5_CTX *context);
	void MD5Update(MD5_CTX *context, const void* input, unsigned int inputLen);
	void MD5UpdaterString(MD5_CTX* context, const char* string);
	int MD5FileUpdateFile(MD5_CTX* context, const char* filename);
	void MD5Final(unsigned char digest[16], MD5_CTX* context);
	void MDString(char *string, unsigned char digest[16]);
	int MD5File(const char *filename, unsigned char digest[16]);
	long long MD5FileExt(const char* filename, unsigned char digest[16]);

#ifdef __cplusplus
}
#endif

#endif /*__MD5C_H__*/
