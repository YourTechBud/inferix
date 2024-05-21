# Use a multi-stage build to keep the final image small
FROM rust:1.77-slim-buster as build

# Install all the OS level dependencies that we have got
RUN apt-get update && apt-get install pkg-config libssl-dev -y && rm -rf /var/lib/apt/lists/*

# Create a new empty shell project
RUN cargo new --bin inferix
WORKDIR /inferix

# Copy over your manifests
COPY ./Cargo.lock ./Cargo.lock
COPY ./Cargo.toml ./Cargo.toml

# Cache dependencies - this layer will be rebuilt only if your dependencies change
RUN cargo build --release --locked
RUN rm src/*.rs

# Copy your source tree
COPY ./src ./src

# Build for release, without the "dev" dependencies
RUN cargo build --release --locked

# Start a new stage from a slim base to reduce image size
FROM debian:buster-slim

# Create a new user and group with a non-root user
RUN groupadd -r inferixuser && useradd -r -g inferixuser inferixuser

# Install only the runtime dependencies
RUN apt-get update && apt-get install -y libssl-dev ca-certificates && rm -rf /var/lib/apt/lists/*

# Switch the working directory to something more familiar
WORKDIR /inferix

# Copy the build artifact from the build stage
COPY --from=build /inferix/target/release/inferix .

# Change the ownership of the binary to the non-root user
RUN chown inferixuser:inferixuser ./inferix

# Switch to the non-root user
USER inferixuser

# Set the startup command to run your binary
CMD ["/inferix/inferix", "--config", "/config/config.yaml"]
