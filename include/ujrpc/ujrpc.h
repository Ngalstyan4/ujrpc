/**
 * @file ujrpc.h
 * @author Ashot Vardanian
 * @date Feb 3, 2023
 * @addtogroup C
 *
 * @brief Binary Interface for Uninterrupted JSON RPC.
 */

#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h> // `bool`
#include <stddef.h>  // `size_t`
#include <stdint.h>  // `int64_t`

typedef void* ujrpc_server_t;
typedef void* ujrpc_call_t;
typedef void* ujrpc_callback_tag_t;
typedef char const* ujrpc_str_t;

typedef void (*ujrpc_callback_t)(ujrpc_call_t, ujrpc_callback_tag_t);

typedef struct ujrpc_config_t {
    char const* interface;
    uint16_t port;
    uint16_t queue_depth;
    uint16_t max_callbacks;
    uint16_t max_threads;

    /// @brief Common choices, aside from a TCP socket are:
    /// > STDOUT_FILENO: console output.
    /// > STDERR_FILENO: errors.
    int32_t logs_file_descriptor;
    /// @brief Can be:
    /// > "human" will print human-readable unit-normalized lines.
    /// > "json" will output newline-delimited JSONs documents.
    char const* logs_format;

    uint16_t max_batch_size;
    uint32_t max_concurrent_connections;
    uint32_t max_lifetime_micro_seconds;
    uint32_t max_lifetime_exchanges;
} ujrpc_config_t;

void ujrpc_init(ujrpc_config_t*, ujrpc_server_t*);
void ujrpc_add_procedure(ujrpc_server_t, ujrpc_str_t, ujrpc_callback_t, ujrpc_callback_tag_t);
void ujrpc_take_call(ujrpc_server_t, uint16_t thread_idx);
void ujrpc_take_calls(ujrpc_server_t, uint16_t thread_idx);
void ujrpc_free(ujrpc_server_t);

bool ujrpc_param_named_bool(ujrpc_call_t, ujrpc_str_t, size_t, bool*);
bool ujrpc_param_named_i64(ujrpc_call_t, ujrpc_str_t, size_t, int64_t*);
bool ujrpc_param_named_f64(ujrpc_call_t, ujrpc_str_t, size_t, double*);
bool ujrpc_param_named_str(ujrpc_call_t, ujrpc_str_t, size_t, ujrpc_str_t*, size_t*);

bool ujrpc_param_positional_bool(ujrpc_call_t, size_t, bool*);
bool ujrpc_param_positional_i64(ujrpc_call_t, size_t, int64_t*);
bool ujrpc_param_positional_f64(ujrpc_call_t, size_t, double*);
bool ujrpc_param_positional_str(ujrpc_call_t, size_t, ujrpc_str_t*, size_t*);

void ujrpc_call_reply_content(ujrpc_call_t, ujrpc_str_t, size_t);
void ujrpc_call_reply_error(ujrpc_call_t, int, ujrpc_str_t, size_t);
void ujrpc_call_reply_error_invalid_params(ujrpc_call_t);
void ujrpc_call_reply_error_out_of_memory(ujrpc_call_t);
void ujrpc_call_reply_error_unknown(ujrpc_call_t);

bool ujrpc_param_named_json(ujrpc_call_t, ujrpc_str_t, size_t, ujrpc_str_t*, size_t*); // TODO
bool ujrpc_param_positional_json(ujrpc_call_t, size_t, ujrpc_str_t*, size_t*);         // TODO

#ifdef __cplusplus
} /* end extern "C" */
#endif