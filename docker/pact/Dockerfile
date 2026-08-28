# Use an official Go image
FROM golang:1.26@sha256:3aff6657219a4d9c14e27fb1d8976c49c29fddb70ba835014f477e1c70636647

# Install dependencies
RUN apt-get update && apt-get install -y \
    curl \
    unzip \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

# Install pact-ffi
RUN curl -LO https://github.com/pact-foundation/pact-ffi/releases/latest/download/libpact_ffi-linux-x86_64.tar.gz \
    && tar -xzf libpact_ffi-linux-x86_64.tar.gz \
    && mv libpact_ffi.so /usr/local/lib/ \
    && ldconfig

# Set up working directory
WORKDIR /app

# Run tests
CMD ["go", "test", "-v", "./..."]
