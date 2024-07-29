#
# Copyright (C) 2024 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

import sqlite3
import os

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

class ModuleNotFoundError(PortError):
    """Exception raised when a module is not found for deallocation."""
    def __init__(self, module_name, message=None):
        self.module_name = module_name
        if message is None:
            message = f"Module '{module_name}' not found."
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

def allocate_ports(required_ports: int, module_name: str, protocol: str):
    """
    Allocate a range of ports for a given module,
    if it is already allocated it is deallocated first.

    :param required_ports: Number of consecutive ports required.
    :param module_name: Name of the module requesting the ports.
    :param protocol: Protocol type ('tcp' or 'udp').
    :return: A tuple (start_port, end_port) if allocation is successful, None otherwise.
    """
    if required_ports < 1:
        raise InvalidPortRequestError()  # Raise error if requested ports are less than 1

    range_start = 20000
    range_end = 45000

    try:
        with sqlite3.connect('./ports.sqlite') as database:
            cursor = database.cursor()
            create_tables(cursor)  # Ensure the tables exist

            # Fetch used ports based on protocol
            if protocol == 'tcp':
                cursor.execute("SELECT * FROM TCP_PORTS ORDER BY start;")
            elif protocol == 'udp':
                cursor.execute("SELECT * FROM UDP_PORTS ORDER BY start;")
            ports_used = cursor.fetchall()

            # If the module already has an assigned range, deallocate it first
            if any(module_name == range[2] for range in ports_used):
                deallocate_ports(module_name, protocol)
                # Reload the used ports after deallocation
                if protocol == 'tcp':
                    cursor.execute("SELECT * FROM TCP_PORTS ORDER BY start;")
                elif protocol == 'udp':
                    cursor.execute("SELECT * FROM UDP_PORTS ORDER BY start;")
                ports_used = cursor.fetchall()

            range_search = True

            while range_search:
                for index in range(required_ports):
                    current_port = range_start + index

                    # Raise an error if the range exceeds the maximum allowed
                    if current_port >= range_end:
                        raise PortRangeExceededError()

                    # Check if the current port is within an already used range
                    for port in ports_used:
                        if current_port in range(port[0], port[1] + 1):
                            range_start = port[1] + 1  # Move to the next available port range
                            break
                else:
                    range_search = False  # Stop searching when a valid range is found

            # Insert the allocated port range into the appropriate table based on protocol
            if protocol == 'tcp':
                cursor.execute("INSERT INTO TCP_PORTS (start, end, module) VALUES (?, ?, ?);", 
                           (range_start, range_start + required_ports - 1, module_name))
            elif protocol == 'udp':
                cursor.execute("INSERT INTO UDP_PORTS (start, end, module) VALUES (?, ?, ?);", 
                           (range_start, range_start + required_ports - 1, module_name))
            database.commit()

            return (range_start, range_start + required_ports - 1)

    except sqlite3.Error as e:
        raise StorageError(f"Database error: {e}")  # Raise custom database error

def deallocate_ports(module_name: str, protocol: str):
    """
    Deallocate the ports for a given module.

    :param module_name: Name of the module whose ports are to be deallocated.
    :param protocol: Protocol type ('tcp' or 'udp').
    :return: A tuple (start_port, end_port) if deallocation is successful, None otherwise.
    """
    try:
        with sqlite3.connect('./ports.sqlite') as database:
            cursor = database.cursor()
            create_tables(cursor)  # Ensure the tables exist

            # Fetch the port range for the given module and protocol
            if protocol == 'tcp':
                cursor.execute("SELECT * FROM TCP_PORTS WHERE module=?;", (module_name,))
            elif protocol == 'udp':
                cursor.execute("SELECT * FROM UDP_PORTS WHERE module=?;", (module_name,))
            ports_deallocated = cursor.fetchall()
            
            if ports_deallocated:
                # Delete the allocated port range for the module
                if protocol == 'tcp':
                    cursor.execute("DELETE FROM TCP_PORTS WHERE module=?;", (module_name,))
                elif protocol == 'udp':
                    cursor.execute("DELETE FROM UDP_PORTS WHERE module=?;", (module_name,))
                database.commit()
                return (ports_deallocated[0][0], ports_deallocated[0][1])
            else:
                raise ModuleNotFoundError(module_name)  # Raise error if the module is not found

    except sqlite3.Error as e:
        raise StorageError(f"Database error: {e}")  # Raise custom database error

def write_range(start: int, end: int, module: str, protocol: str):
    """
    Write a port range for a module directly to the database.

    :param start: Starting port number.
    :param end: Ending port number.
    :param module: Name of the module.
    :param protocol: Protocol type ('tcp' or 'udp').
    """
    try:
        with sqlite3.connect('./ports.sqlite') as database:
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
        raise StorageError(f"Database error: {e}")  # Raise custom database error
