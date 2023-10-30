#ifndef __DUMP1090_H__
#define __DUMP1090_H__

#include <stdint.h>
#include <stdlib.h>
#include <unistd.h>



#define MODES_LONG_MSG_BYTES (112/8)
#define MODES_SHORT_MSG_BYTES (56/8)



int startProcess(uint32_t deviceIndex, int gain, uint32_t frequency, uint8_t enableAGC, char* filename, int loop);

/* The struct we use to store information about a decoded message. */
typedef struct {
    /* Generic fields */
    unsigned char msg[MODES_LONG_MSG_BYTES]; /* Binary message. */
    int msgbits;                /* Number of bits in message */
    int msgtype;                /* Downlink format # */
    int crcok;                  /* True if CRC was valid */
    uint32_t crc;               /* Message CRC */
    int errorbit;               /* Bit corrected. -1 if no bit corrected. */
    int aa1, aa2, aa3;          /* ICAO Address bytes 1 2 and 3 */
    int phase_corrected;        /* True if phase correction was applied. */

    /* DF 11 */
    int ca;                     /* Responder capabilities. */

    /* DF 17 */
    int metype;                 /* Extended squitter message type. */
    int mesub;                  /* Extended squitter message subtype. */
    int heading_is_valid;
    int heading;
    int aircraft_type;
    int fflag;                  /* 1 = Odd, 0 = Even CPR message. */
    int tflag;                  /* UTC synchronized? */
    int raw_latitude;           /* Non decoded latitude */
    int raw_longitude;          /* Non decoded longitude */
    char flight[9];             /* 8 chars flight number. */
    int ew_dir;                 /* 0 = East, 1 = West. */
    int ew_velocity;            /* E/W velocity. */
    int ns_dir;                 /* 0 = North, 1 = South. */
    int ns_velocity;            /* N/S velocity. */
    int vert_rate_source;       /* Vertical rate source. */
    int vert_rate_sign;         /* Vertical rate sign. */
    int vert_rate;              /* Vertical rate. */
    int velocity;               /* Computed from EW and NS velocity. */

    /* DF4, DF5, DF20, DF21 */
    int fs;                     /* Flight status for DF4,5,20,21 */
    int dr;                     /* Request extraction of downlink request. */
    int um;                     /* Request extraction of downlink request. */
    int identity;               /* 13 bits identity (Squawk). */

    /* Fields used by multiple message types. */
    int altitude, unit;
} modesMessage;

/* Structure used to describe an aircraft in iteractive mode. */
typedef struct Aircraft{
    uint32_t addr;      /* ICAO address */
    char hexaddr[7];    /* Printable ICAO address */                           /*`json:"hex"`*/
    char flight[9];     /* Flight number */                                    /*`json:"flight"`*/
    int altitude;       /* Altitude */                                         /*`json:"altitude"`*/
    int speed;          /* Velocity computed from EW and NS components. */     /*`json:"speed"`*/
    int track;          /* Angle of flight. */                                 /*`json:"track"`*/
    time_t seen;        /* Time at which the last packet was received. */
    long messages;      /* Number of Mode S messages received. */
    /* Encoded latitude and longitude as extracted by odd and even
     * CPR encoded messages. */
    int odd_cprlat;
    int odd_cprlon;
    int even_cprlat;
    int even_cprlon;
    double lat, lon;    /* Coordinated obtained from CPR encoded data. */    /*`json:"lat"`*/  /*`json:"lon"`*/
    long long odd_cprtime, even_cprtime;
    struct Aircraft *next; /* Next aircraft in our linked list. */
} aircraft;

extern void goSendMessage(modesMessage*);
extern void goSendAircraft(modesMessage*, aircraft*);

#endif