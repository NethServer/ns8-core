#
# Copyright (C) 2024 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

import sqlite3

class PortError(Exception):
    """Base class for all port-related exceptions."""
    pass

class PortRangeExceededError(PortError):
    """Exception raised when the port range is exceeded."""
    def __init__(self, message="Ports range max exceeded!"):
        self.message = message
        super().__init__(self.message)

class StorageError(PortError):
    """Exception raised when a database error occurs."""
    def __init__(self, message="Database operation failed."):
        self.message = message
        super().__init__(self.message)

class InvalidPortRequestError(PortError):
    """Exception raised when the requested number of ports is invalid."""
    def __init__(self, message="The number of required ports must be at least 1."):
        self.message = message
        super().__init__(self.message)

def create_tables(cursor: sqlite3.Cursor):
    # Create TCP table if it doesn't exist
    cursor.execute("""
        CREATE TABLE IF NOT EXISTS TCP_PORTS (
            start INT NOT NULL,
            end INT NOT NULL,
            module CHAR(255) NOT NULL
        );
    """)

    # Create UDP table if it doesn't exist
    cursor.execute("""
        CREATE TABLE IF NOT EXISTS UDP_PORTS (
            start INT NOT NULL,
            end INT NOT NULL,
            module CHAR(255) NOT NULL
        );
    """)

def is_port_used(ports_used, port_to_check):
    for port in ports_used:
        if port_to_check in range(port[0], port[1] + 1):
            return True
    return False

def allocate_ports(required_ports: int, module_name: str, protocol: str, keep_existing: bool=False):
    """
    Allocate a range of ports for a given module,
    if it is already allocated it is deallocated first.

    :param required_ports: Number of consecutive ports required.
    :param module_name: Name of the module requesting the ports.
    :param protocol: Protocol type ('tcp' or 'udp').
    :param keep_existing: If True, keep the existing port range for the module.
    :return: A tuple (start_port, end_port) if allocation is successful, None otherwise.
    """
    if required_ports < 1:
        raise InvalidPortRequestError()  # Raise error if requested ports are less than 1

    range_start = 20000
    range_end = 45000

    try:
        with sqlite3.connect('./ports.sqlite', isolation_level='EXCLUSIVE', timeout=30) as database:
            cursor = database.cursor()
            create_tables(cursor) # Ensure the tables exist

            # Fetch used ports based on protocol
            if protocol == 'tcp':
                cursor.execute("SELECT start,end,module FROM TCP_PORTS ORDER BY start;")
            elif protocol == 'udp':
                cursor.execute("SELECT start,end,module FROM UDP_PORTS ORDER BY start;")
            ports_used = cursor.fetchall()

            # If the module already has an assigned range, deallocate it first
            if any(module_name == range[2] for range in ports_used):
                if not keep_existing:
                    deallocate_ports(module_name, protocol)
                # Reload the used ports after deallocation
                if protocol == 'tcp':
                    cursor.execute("SELECT start,end,module FROM TCP_PORTS ORDER BY start;")
                elif protocol == 'udp':
                    cursor.execute("SELECT start,end,module FROM UDP_PORTS ORDER BY start;")
                ports_used = cursor.fetchall()

            if len(ports_used) == 0:
                write_range(range_start, range_start + required_ports - 1, module_name, protocol, database)
                return (range_start, range_start + required_ports - 1)

            while range_start <= range_end:
                # Check if the current port is within an already used range
                for port_range in ports_used:
                    for index in range(required_ports):
                        if is_port_used(ports_used, range_start+index):
                            range_start = port_range[1] + 1  # Move to the next available port range
                            break
                        if index == required_ports-1:
                            write_range(range_start, range_start + required_ports - 1, module_name, protocol, database)
                            return (range_start, range_start + required_ports - 1)
            else:
                raise PortRangeExceededError()
    except sqlite3.Error as e:
        raise StorageError(f"Database error: {e}") from e # Raise custom database error

def deallocate_ports(module_name: str, protocol: str):
    """
    Deallocate the ports for a given module.

    :param module_name: Name of the module whose ports are to be deallocated.
    :param protocol: Protocol type ('tcp' or 'udp').
    :return: a list of tuples if deallocation is successful, None otherwise.
    """
    try:
        with sqlite3.connect('./ports.sqlite', isolation_level='EXCLUSIVE', timeout=30) as database:
            cursor = database.cursor()
            create_tables(cursor)  # Ensure the tables exist

            # Fetch the port range for the given module and protocol
            if protocol == 'tcp':
                cursor.execute("SELECT start,end,module FROM TCP_PORTS WHERE module=?;", (module_name,))
            elif protocol == 'udp':
                cursor.execute("SELECT start,end,module FROM UDP_PORTS WHERE module=?;", (module_name,))
            ports_deallocated = cursor.fetchall()

            if ports_deallocated:
                # Delete the allocated port range for the module
                if protocol == 'tcp':
                    cursor.execute("DELETE FROM TCP_PORTS WHERE module=?;", (module_name,))
                elif protocol == 'udp':
                    cursor.execute("DELETE FROM UDP_PORTS WHERE module=?;", (module_name,))
                database.commit()
                # remove the name of the module and return a list of tuples
                return [(port[0], port[1]) for port in ports_deallocated]
            else:
                return []

    except sqlite3.Error as e:
        raise StorageError(f"Database error: {e}") from e # Raise custom database error

def write_range(start: int, end: int, module: str, protocol: str, database: sqlite3.Connection=None):
    """
    Write a port range for a module directly to the database.

    :param start: Starting port number.
    :param end: Ending port number.
    :param module: Name of the module.
    :param protocol: Protocol type ('tcp' or 'udp').
    """
    try:
        if database is None:
            database = sqlite3.connect('./ports.sqlite', isolation_level='EXCLUSIVE', timeout=30)

        with database:
            cursor = database.cursor()
            create_tables(cursor)  # Ensure the tables exist

            # Insert the port range into the appropriate table based on protocol
            if protocol == 'tcp':
                cursor.execute("INSERT INTO TCP_PORTS (start, end, module) VALUES (?, ?, ?);", 
                           (start, end, module))
            elif protocol == 'udp':
                cursor.execute("INSERT INTO UDP_PORTS (start, end, module) VALUES (?, ?, ?);", 
                           (start, end, module))
            database.commit()

    except sqlite3.Error as e:
        raise StorageError(f"Database error: {e}") from e # Raise custom database error

def get_tcp_ports_by_module(module_name: str) -> int:
    """
    Get the number of TCP ports allocated to a specific module.

    :param module_name: Name of the module.
    :return: Number of TCP ports allocated to the module.
    """
    try:
        with sqlite3.connect('./ports.sqlite', isolation_level='EXCLUSIVE', timeout=30) as database:
            cursor = database.cursor()
            create_tables(cursor)

            cursor.execute("""
                SELECT SUM(end - start + 1) 
                FROM TCP_PORTS 
                WHERE module=?;
            """, (module_name,))
            result = cursor.fetchone()

            return result[0] if result[0] is not None else 0

    except sqlite3.Error as e:
        raise StorageError(f"Database error: {e}") from e

def get_udp_ports_by_module(module_name: str) -> int:
    """
    Get the number of UDP ports allocated to a specific module.

    :param module_name: Name of the module.
    :return: Number of UDP ports allocated to the module.
    """
    try:
        with sqlite3.connect('./ports.sqlite', isolation_level='EXCLUSIVE', timeout=30) as database:
            cursor = database.cursor()
            create_tables(cursor)

            cursor.execute("""
                SELECT SUM(end - start + 1) 
                FROM UDP_PORTS 
                WHERE module=?;
            """, (module_name,))
            result = cursor.fetchone()

            return result[0] if result[0] is not None else 0

    except sqlite3.Error as e:
        raise StorageError(f"Database error: {e}") from e
